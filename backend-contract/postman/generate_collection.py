#!/usr/bin/env python3
"""从冻结 OpenAPI 与权威成功 fixture 生成 Postman Collection v2.1。"""

from __future__ import annotations

import argparse
from copy import deepcopy
import hashlib
from http import HTTPStatus
import json
from pathlib import Path
import re
import sys
import uuid


CONTRACT_ROOT = Path(__file__).resolve().parents[1]
OPENAPI_PATH = CONTRACT_ROOT / "openapi" / "clearwave-backend-v1.json"
EXAMPLES_DIR = CONTRACT_ROOT / "examples"
OUTPUT_PATH = Path(__file__).with_name(
    "Clearwave-Backend-v1.postman_collection.json"
)
USAGE_GUIDE_PATH = Path(__file__).with_name("API_FUNCTION_USAGE_GUIDE.md")
POSTMAN_SCHEMA = (
    "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
)
HTTP_METHODS = ("get", "post", "put", "patch", "delete", "head", "options")
MUTATION_METHODS = frozenset({"POST", "PUT", "PATCH", "DELETE"})

SENSITIVE_BODY_VARIABLES = {
    "google_id_token": "{{googleIdToken}}",
    "guest_upgrade_grant": "{{guestUpgradeGrant}}",
    "invite_code": "{{inviteCode}}",
    "purchase_token": "{{purchaseToken}}",
}
DYNAMIC_NONCES = frozenset({"client_nonce", "proof_nonce", "redemption_nonce"})
AUTH_VARIABLE_BY_OPERATION = {
    "exchangeGoogleProof": "inviteAccessToken",
    "redeemAdministratorInvite": "inviteAccessToken",
    "redeemMemberInvite": "inviteAccessToken",
    "refreshSession": "refreshToken",
}


def load_json(path: Path) -> dict[str, object]:
    """读取 JSON 对象，并在根节点不是对象时拒绝继续生成。"""

    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"JSON root must be an object: {path}")
    return value


def resolve_ref(document: dict[str, object], value: object) -> object:
    """解析当前 OpenAPI 文档内部的 JSON Pointer，拒绝外部引用。"""

    resolved = value
    seen: set[str] = set()
    while isinstance(resolved, dict) and "$ref" in resolved:
        reference = resolved["$ref"]
        if not isinstance(reference, str) or not reference.startswith("#/"):
            raise ValueError(f"unsupported external reference: {reference!r}")
        if reference in seen:
            raise ValueError(f"cyclic reference: {reference}")
        seen.add(reference)
        target: object = document
        for raw_part in reference[2:].split("/"):
            part = raw_part.replace("~1", "/").replace("~0", "~")
            if not isinstance(target, dict) or part not in target:
                raise ValueError(f"unresolved reference: {reference}")
            target = target[part]
        resolved = target
    return resolved


def collect_operations(
    document: dict[str, object],
) -> dict[str, tuple[str, str, dict[str, object], dict[str, object]]]:
    """按 operationId 收集 method、path、operation 与 path item。"""

    paths = document.get("paths")
    if not isinstance(paths, dict):
        raise ValueError("OpenAPI paths must be an object")
    operations: dict[
        str, tuple[str, str, dict[str, object], dict[str, object]]
    ] = {}
    for path, raw_path_item in paths.items():
        if not isinstance(path, str) or not isinstance(raw_path_item, dict):
            raise ValueError("invalid OpenAPI path item")
        path_item = raw_path_item
        for method in HTTP_METHODS:
            raw_operation = path_item.get(method)
            if raw_operation is None:
                continue
            if not isinstance(raw_operation, dict):
                raise ValueError(f"invalid operation: {method.upper()} {path}")
            operation_id = raw_operation.get("operationId")
            if not isinstance(operation_id, str) or not operation_id:
                raise ValueError(f"missing operationId: {method.upper()} {path}")
            if operation_id in operations:
                raise ValueError(f"duplicate operationId: {operation_id}")
            operations[operation_id] = (
                method.upper(),
                path,
                raw_operation,
                path_item,
            )
    return operations


def collect_success_fixtures() -> tuple[dict[str, dict[str, object]], str]:
    """为每个 operation 选取首次提交或普通成功 fixture，并计算 fixture 摘要。"""

    candidates: dict[str, list[dict[str, object]]] = {}
    digest = hashlib.sha256()
    for path in sorted(EXAMPLES_DIR.glob("*.json")):
        raw = path.read_bytes()
        digest.update(path.name.encode("utf-8"))
        digest.update(b"\0")
        digest.update(raw)
        document = json.loads(raw)
        fixtures = document.get("fixtures") if isinstance(document, dict) else None
        if not isinstance(fixtures, list):
            raise ValueError(f"fixtures must be an array: {path}")
        for fixture in fixtures:
            if not isinstance(fixture, dict) or fixture.get("kind") != "success":
                continue
            operation_id = fixture.get("operation_id")
            if not isinstance(operation_id, str) or not operation_id:
                raise ValueError(f"success fixture has no operation_id: {path}")
            candidates.setdefault(operation_id, []).append(fixture)

    def priority(fixture: dict[str, object]) -> tuple[int, str]:
        case_id = str(fixture.get("case_id") or "")
        if case_id == "success_first_committed":
            return 0, case_id
        if case_id == "success":
            return 1, case_id
        if "replay" in case_id:
            return 3, case_id
        return 2, case_id

    selected = {
        operation_id: sorted(fixtures, key=priority)[0]
        for operation_id, fixtures in candidates.items()
    }
    return selected, digest.hexdigest()


def operation_parameters(
    document: dict[str, object],
    path_item: dict[str, object],
    operation: dict[str, object],
) -> list[dict[str, object]]:
    """合并 path-level 与 operation-level 参数，后者覆盖同名参数。"""

    merged: dict[tuple[str, str], dict[str, object]] = {}
    for owner in (path_item, operation):
        raw_parameters = owner.get("parameters", [])
        if not isinstance(raw_parameters, list):
            raise ValueError("operation parameters must be an array")
        for raw_parameter in raw_parameters:
            parameter = resolve_ref(document, raw_parameter)
            if not isinstance(parameter, dict):
                raise ValueError("resolved parameter must be an object")
            name = parameter.get("name")
            location = parameter.get("in")
            if not isinstance(name, str) or not isinstance(location, str):
                raise ValueError("parameter must have name and in")
            merged[(location, name)] = parameter
    return list(merged.values())


def snake_to_camel(value: str) -> str:
    """把 wire snake_case 名转换为稳定的 Postman 变量名。"""

    head, *tail = value.split("_")
    return head + "".join(part[:1].upper() + part[1:] for part in tail)


def scalar_text(value: object) -> str:
    """把 fixture 标量转换为 Postman 参数文本。"""

    if value is None:
        return ""
    if isinstance(value, bool):
        return "true" if value else "false"
    if isinstance(value, (dict, list)):
        return json.dumps(value, ensure_ascii=False, separators=(",", ":"))
    return str(value)


def parameter_example(
    document: dict[str, object], parameter: dict[str, object]
) -> object:
    """为 fixture 未启用的可选参数选择受合同约束的展示值。"""

    if "example" in parameter:
        return parameter["example"]
    schema = resolve_ref(document, parameter.get("schema", {}))
    if not isinstance(schema, dict):
        return ""
    for key in ("example", "default", "const"):
        if key in schema:
            return schema[key]
    enum = schema.get("enum")
    if isinstance(enum, list) and enum:
        return enum[0]
    return ""


def substitute_sensitive_body(value: object, key: str | None = None) -> object:
    """把敏感 fixture 输入替换为未赋值变量，并为 nonce 使用动态 GUID。"""

    if key in SENSITIVE_BODY_VARIABLES:
        return SENSITIVE_BODY_VARIABLES[key]
    if key in DYNAMIC_NONCES:
        return "{{$guid}}"
    if isinstance(value, dict):
        return {
            child_key: substitute_sensitive_body(child_value, child_key)
            for child_key, child_value in value.items()
        }
    if isinstance(value, list):
        return [substitute_sensitive_body(item) for item in value]
    return value


def fixture_request(fixture: dict[str, object]) -> dict[str, object]:
    """返回 fixture 请求对象并验证其基本形状。"""

    request = fixture.get("request")
    if not isinstance(request, dict):
        raise ValueError(f"fixture request must be an object: {fixture.get('operation_id')}")
    for key in ("headers_present", "path_parameters", "query_parameters"):
        value = request.get(key, [] if key == "headers_present" else {})
        if key == "headers_present" and not isinstance(value, list):
            raise ValueError(f"{key} must be an array")
        if key != "headers_present" and not isinstance(value, dict):
            raise ValueError(f"{key} must be an object")
    return request


def request_auth(
    operation_id: str, request_fixture: dict[str, object]
) -> dict[str, object]:
    """按 fixture 的 Authorization 要求选择对应凭据变量。"""

    headers = request_fixture.get("headers_present", [])
    if "Authorization" not in headers:
        return {"type": "noauth"}
    variable = AUTH_VARIABLE_BY_OPERATION.get(operation_id, "accessToken")
    return {
        "type": "bearer",
        "bearer": [{"key": "token", "value": f"{{{{{variable}}}}}", "type": "string"}],
    }


def request_headers(request_fixture: dict[str, object]) -> list[dict[str, object]]:
    """从权威 fixture 构造请求头；Authorization 交给 Postman auth owner。"""

    values = {
        "X-Request-Id": "{{$guid}}",
        "X-Contract-Version": "1",
        "X-Client-Version": "{{clientVersion}}",
        "Idempotency-Key": "{{idempotencyKey}}",
    }
    headers: list[dict[str, object]] = []
    for raw_name in request_fixture.get("headers_present", []):
        name = str(raw_name)
        if name == "Authorization":
            continue
        headers.append(
            {
                "key": name,
                "value": values.get(name, ""),
                "type": "text",
            }
        )
    if "body" in request_fixture:
        headers.append({"key": "Content-Type", "value": "application/json", "type": "text"})
    return headers


def build_url(
    document: dict[str, object],
    path: str,
    parameters: list[dict[str, object]],
    request_fixture: dict[str, object],
) -> tuple[dict[str, object], dict[str, str]]:
    """构造 Postman URL，并返回 path 变量的合成默认值。"""

    raw_path_parameters = request_fixture.get("path_parameters", {})
    raw_query_parameters = request_fixture.get("query_parameters", {})
    assert isinstance(raw_path_parameters, dict)
    assert isinstance(raw_query_parameters, dict)

    path_defaults: dict[str, str] = {}

    def replace_path(match: re.Match[str]) -> str:
        wire_name = match.group(1)
        if wire_name not in raw_path_parameters:
            raise ValueError(f"fixture misses path parameter {wire_name}: {path}")
        variable_name = snake_to_camel(wire_name)
        path_defaults[variable_name] = scalar_text(raw_path_parameters[wire_name])
        return f"{{{{{variable_name}}}}}"

    postman_path = re.sub(r"\{([^{}]+)\}", replace_path, path)
    known_query_names = {
        str(parameter.get("name"))
        for parameter in parameters
        if parameter.get("in") == "query"
    }
    unknown_query = set(raw_query_parameters) - known_query_names
    if unknown_query:
        raise ValueError(f"fixture contains unknown query parameters: {sorted(unknown_query)}")

    query: list[dict[str, object]] = []
    for parameter in parameters:
        if parameter.get("in") != "query":
            continue
        name = str(parameter["name"])
        enabled = name in raw_query_parameters
        if parameter.get("required") is True and not enabled:
            raise ValueError(f"fixture misses required query parameter {name}: {path}")
        value = (
            raw_query_parameters[name]
            if enabled
            else parameter_example(document, parameter)
        )
        item: dict[str, object] = {
            "key": name,
            "value": scalar_text(value),
            "description": str(parameter.get("description") or ""),
        }
        if not enabled:
            item["disabled"] = True
        query.append(item)

    enabled_query = [item for item in query if item.get("disabled") is not True]
    raw_query = "&".join(f"{item['key']}={item['value']}" for item in enabled_query)
    raw_url = "{{baseUrl}}" + postman_path
    if raw_query:
        raw_url += "?" + raw_query
    return (
        {
            "raw": raw_url,
            "host": ["{{baseUrl}}"],
            "path": postman_path.lstrip("/").split("/"),
            "query": query,
        },
        path_defaults,
    )


def build_request(
    document: dict[str, object],
    operation_id: str,
    method: str,
    path: str,
    operation: dict[str, object],
    path_item: dict[str, object],
    fixture: dict[str, object],
) -> tuple[dict[str, object], dict[str, str]]:
    """构建单个 Postman request，并携带选定成功 fixture 的来源说明。"""

    request_fixture = fixture_request(fixture)
    parameters = operation_parameters(document, path_item, operation)
    url, path_defaults = build_url(document, path, parameters, request_fixture)
    summary = str(operation.get("summary") or operation_id)
    description = str(operation.get("description") or "")
    fixture_id = str(fixture.get("case_id") or "")
    precondition = str(fixture.get("precondition") or "")
    request: dict[str, object] = {
        "auth": request_auth(operation_id, request_fixture),
        "method": method,
        "header": request_headers(request_fixture),
        "url": url,
        "description": (
            f"### `{operation_id}`\n\n"
            f"`{method} {path}`\n\n"
            f"{description}\n\n"
            f"示例来源：`{fixture_id}`；前置状态：`{precondition}`。\n\n"
            "请求示例只使用合成 fixture。完整成功/错误响应以 "
            "`backend-contract/examples/*.json` 为准。"
        ),
    }
    if "body" in request_fixture:
        body = substitute_sensitive_body(request_fixture["body"])
        request["body"] = {
            "mode": "raw",
            "raw": json.dumps(body, ensure_ascii=False, indent=2),
            "options": {"raw": {"language": "json"}},
        }
    return request, path_defaults


def build_response(
    request: dict[str, object], fixture: dict[str, object]
) -> dict[str, object]:
    """把选定成功 fixture 附为 Postman 示例响应。"""

    response = fixture.get("response")
    if not isinstance(response, dict):
        raise ValueError(f"fixture response must be an object: {fixture.get('operation_id')}")
    code = response.get("status")
    if not isinstance(code, int):
        raise ValueError(f"fixture response status must be an integer: {fixture.get('operation_id')}")
    body = response.get("body")
    try:
        status = HTTPStatus(code).phrase
    except ValueError:
        status = str(code)
    return {
        "name": f"{fixture.get('case_id')} ({code})",
        "originalRequest": deepcopy(request),
        "status": status,
        "code": code,
        "_postman_previewlanguage": "json",
        "header": [{"key": "Content-Type", "value": "application/json"}],
        "cookie": [],
        "body": json.dumps(body, ensure_ascii=False, indent=2),
    }


def build_test_event(expected_status: int) -> dict[str, object]:
    """为每个请求添加稳定成功 envelope 的轻量断言。"""

    lines = [
        f'pm.test("HTTP {expected_status}", function () {{',
        f"    pm.response.to.have.status({expected_status});",
        "});",
        'pm.test("Clearwave contract envelope", function () {',
        "    const payload = pm.response.json();",
        '    pm.expect(payload.contract_version).to.eql("1");',
        "    pm.expect(payload.request_id).to.be.a(\"string\").and.not.empty;",
        "    pm.expect(payload.server_time_ms).to.be.a(\"number\");",
        "});",
    ]
    return {
        "listen": "test",
        "script": {"type": "text/javascript", "exec": lines},
    }


def collection_pre_request_event() -> dict[str, object]:
    """按 operation 持久化幂等键，使重复 Send 能验证 canonical replay。"""

    lines = [
        'const mutationMethods = ["POST", "PUT", "PATCH", "DELETE"];',
        "if (mutationMethods.includes(pm.request.method)) {",
        '    const operationId = pm.info.requestName.split(" · ")[0];',
        '    const storageKey = "_idempotencyKey_" + operationId;',
        "    let idempotencyKey = pm.collectionVariables.get(storageKey);",
        "    if (!idempotencyKey) {",
        '        idempotencyKey = pm.variables.replaceIn("{{$guid}}");',
        "        pm.collectionVariables.set(storageKey, idempotencyKey);",
        "    }",
        '    pm.variables.set("idempotencyKey", idempotencyKey);',
        "}",
    ]
    return {
        "listen": "prerequest",
        "script": {"type": "text/javascript", "exec": lines},
    }


def collection_variables(path_defaults: dict[str, str]) -> list[dict[str, object]]:
    """构造无真实凭据的集合变量及稳定合成 ID。"""

    variables: list[tuple[str, str, str]] = [
        (
            "baseUrl",
            "https://staging.invalid",
            "必须替换为合法、非 placeholder、非 loopback 的 HTTPS staging origin。",
        ),
        ("clientVersion", "postman-contract-v1", "X-Client-Version。"),
        ("accessToken", "", "普通 account operation 的 access token；不得提交真实值。"),
        ("refreshToken", "", "refreshSession 的 refresh credential；不得提交真实值。"),
        (
            "inviteAccessToken",
            "",
            "Invite redeem/Guest upgrade 使用的受限 Guest 或 account token。",
        ),
        ("googleIdToken", "", "真实 Google proof；不得提交到仓库或共享集合。"),
        ("guestUpgradeGrant", "", "Guest upgrade grant；与 token 同级保护。"),
        ("inviteCode", "", "一次性邀请码；不得记录或提交。"),
        ("purchaseToken", "", "Google Play purchase token；不得记录或提交。"),
    ]
    variables.extend(
        (name, value, "来自合成 success fixture 的稳定 path 参数。")
        for name, value in sorted(path_defaults.items())
    )
    return [
        {"key": key, "value": value, "type": "string", "description": description}
        for key, value, description in variables
    ]


def build_collection() -> dict[str, object]:
    """生成完整 Collection，并拒绝 operation/fixture/tag 覆盖缺口。"""

    document = load_json(OPENAPI_PATH)
    operations = collect_operations(document)
    fixtures, fixtures_sha256 = collect_success_fixtures()
    if set(fixtures) != set(operations):
        missing = sorted(set(operations) - set(fixtures))
        orphaned = sorted(set(fixtures) - set(operations))
        raise ValueError(
            f"success fixture coverage mismatch: missing={missing}, orphaned={orphaned}"
        )

    openapi_bytes = OPENAPI_PATH.read_bytes()
    openapi_sha256 = hashlib.sha256(openapi_bytes).hexdigest()
    info = document.get("info")
    if not isinstance(info, dict):
        raise ValueError("OpenAPI info must be an object")
    version = str(info.get("version") or "")
    tags = document.get("tags", [])
    if not isinstance(tags, list):
        raise ValueError("OpenAPI tags must be an array")
    tag_order = [
        str(tag["name"])
        for tag in tags
        if isinstance(tag, dict) and isinstance(tag.get("name"), str)
    ]
    folders: dict[str, dict[str, object]] = {
        tag: {
            "name": tag,
            "description": f"OpenAPI Tag: `{tag}`",
            "item": [],
        }
        for tag in tag_order
    }
    path_defaults: dict[str, str] = {}

    for operation_id, (method, path, operation, path_item) in operations.items():
        raw_tags = operation.get("tags")
        if not isinstance(raw_tags, list) or len(raw_tags) != 1:
            raise ValueError(f"operation must have exactly one tag: {operation_id}")
        tag = str(raw_tags[0])
        if tag not in folders:
            raise ValueError(f"operation references undeclared tag: {operation_id} -> {tag}")
        fixture = fixtures[operation_id]
        request, request_path_defaults = build_request(
            document,
            operation_id,
            method,
            path,
            operation,
            path_item,
            fixture,
        )
        for name, value in request_path_defaults.items():
            existing = path_defaults.get(name)
            if existing is not None and existing != value:
                raise ValueError(f"conflicting fixture path variable {name}: {existing} != {value}")
            path_defaults[name] = value
        response = fixture.get("response")
        if not isinstance(response, dict) or not isinstance(response.get("status"), int):
            raise ValueError(f"invalid success response: {operation_id}")
        summary = str(operation.get("summary") or operation_id)
        folders[tag]["item"].append(
            {
                "name": f"{operation_id} · {summary}",
                "event": [build_test_event(int(response["status"]))],
                "request": request,
                "response": [build_response(request, fixture)],
            }
        )

    for tag, folder in folders.items():
        if not folder["item"]:
            raise ValueError(f"declared tag has no operations: {tag}")

    collection_id = uuid.uuid5(
        uuid.NAMESPACE_URL,
        f"clearwave-backend:{version}:{openapi_sha256}:{fixtures_sha256}",
    )
    return {
        "info": {
            "_postman_id": str(collection_id),
            "name": "Clearwave Backend Contract v1",
            "description": (
                f"由冻结 OpenAPI `{version}` 与权威 success fixtures 机械生成。\n\n"
                f"- operation：{len(operations)}\n"
                f"- OpenAPI SHA-256：`{openapi_sha256}`\n"
                f"- Fixtures SHA-256：`{fixtures_sha256}`\n\n"
                "默认 `baseUrl=https://staging.invalid`，所有凭据变量为空；替换前不会连接真实服务。"
            ),
            "schema": POSTMAN_SCHEMA,
            "version": version,
        },
        "event": [collection_pre_request_event()],
        "variable": collection_variables(path_defaults),
        "item": [folders[tag] for tag in tag_order],
    }


def rendered_collection() -> str:
    """以稳定 UTF-8/缩进格式序列化生成结果。"""

    return json.dumps(build_collection(), ensure_ascii=False, indent=2) + "\n"


def validate_usage_guide() -> int:
    """校验用途文档恰好覆盖当前 operation，并保持 method/path 一致。"""

    document = load_json(OPENAPI_PATH)
    operations = collect_operations(document)
    if not USAGE_GUIDE_PATH.is_file():
        raise ValueError(f"missing API usage guide: {USAGE_GUIDE_PATH}")
    # Markdown 会为表格列对齐补水平空格；只放宽同行空格，避免残缺表格跨行误匹配。
    row_pattern = re.compile(
        r"^\|[ \t]*`([A-Za-z][A-Za-z0-9]+)`[ \t]*\|[ \t]*"
        r"`([A-Z]+) ([^`\r\n]+)`[ \t]*\|[ \t]*"
        r"([^|\r\n]+?)[ \t]*\|[ \t]*([^|\r\n]+?)[ \t]*\|[ \t]*"
        r"([^|\r\n]+?)[ \t]*\|[ \t]*$",
        re.MULTILINE,
    )
    documented: dict[str, tuple[str, str]] = {}
    for match in row_pattern.finditer(USAGE_GUIDE_PATH.read_text(encoding="utf-8")):
        operation_id, method, path, purpose, feature, owner = match.groups()
        if operation_id in documented:
            raise ValueError(f"duplicate API usage guide row: {operation_id}")
        if not purpose.strip() or not feature.strip() or not owner.strip():
            raise ValueError(f"incomplete API usage guide row: {operation_id}")
        documented[operation_id] = (method, path)
    expected = {
        operation_id: (method, path)
        for operation_id, (method, path, _, _) in operations.items()
    }
    if documented != expected:
        missing = sorted(set(expected) - set(documented))
        orphaned = sorted(set(documented) - set(expected))
        mismatched = sorted(
            operation_id
            for operation_id in set(expected) & set(documented)
            if expected[operation_id] != documented[operation_id]
        )
        raise ValueError(
            "API usage guide mismatch: "
            f"missing={missing}, orphaned={orphaned}, mismatched={mismatched}"
        )
    return len(documented)


def main() -> int:
    """生成集合，或在 --check 模式检查已跟踪文件是否新鲜。"""

    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--check",
        action="store_true",
        help="仅检查现有 collection 与当前合同是否完全一致",
    )
    args = parser.parse_args()
    expected = rendered_collection()
    if args.check:
        usage_operation_count = validate_usage_guide()
        if not OUTPUT_PATH.is_file():
            print(f"missing generated collection: {OUTPUT_PATH}", file=sys.stderr)
            return 1
        if OUTPUT_PATH.read_text(encoding="utf-8") != expected:
            print(
                "generated Postman collection is stale; run "
                "python3 backend-contract/postman/generate_collection.py",
                file=sys.stderr,
            )
            return 1
        collection = json.loads(expected)
        request_count = sum(len(folder["item"]) for folder in collection["item"])
        print(
            "Postman collection and API usage guide passed: "
            f"{request_count} collection operations, "
            f"{usage_operation_count} documented operations"
        )
        return 0
    OUTPUT_PATH.write_text(expected, encoding="utf-8")
    collection = json.loads(expected)
    request_count = sum(len(folder["item"]) for folder in collection["item"])
    print(f"Generated {OUTPUT_PATH}: {request_count} operations")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
