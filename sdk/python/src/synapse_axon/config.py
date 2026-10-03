"""Fail-fast TOML mapping to the maintained flat discovery contract."""
import copy
import math
import re
import tomllib
import warnings
from pathlib import Path

ID = re.compile(r"^[A-Za-z0-9_-]+$")
META = {"id", "name", "group", "tags", "icon", "url", "ttl", "description", "markdown_docs"}
COMMON = {"type", "label", "default", "unit", "monitors"}
PROPS = {
    "stat": {"copyable"}, "gauge": {"min", "max", "thresholds"},
    "status_indicator": {"mapping"}, "log_stream": {"max_items"},
    "action_group": {"items"}, "link": {"uri", "text"},
}

def require(condition, message):
    if not condition:
        raise ValueError(message)

def identifier(value):
    return isinstance(value, str) and ID.fullmatch(value) is not None

def number(value):
    return type(value) in (int, float) and math.isfinite(value)

def keys(value, allowed, location):
    require(isinstance(value, dict), f"{location} must be a table")
    require(not set(value) - allowed, f"unsupported fields in {location}")

def validate_value(component, value):
    kind = component["type"]
    if kind == "gauge":
        require(number(value), "gauge value must be a finite number")
        require(component["min"] <= value <= component["max"], "gauge value outside min/max")
    elif kind == "log_stream":
        require(isinstance(value, str) or (isinstance(value, list) and all(isinstance(x, str) for x in value)), "log_stream value must be a string or string array")
    elif kind == "status_indicator":
        require(isinstance(value,str) and value in component["mapping"], "status_indicator value must have a mapping")
    elif kind == "stat":
        require(value is None or type(value) in (str, bool, int, float), f"{kind} value must be scalar")
        require(not isinstance(value, float) or math.isfinite(value), "value must be finite")
    else:
        raise ValueError(f"{kind} has no mutable value")

def load_config(path):
    with Path(path).open("rb") as stream:
        config = tomllib.load(stream)
    keys(config, {"schema", "meta", "layout", "components"}, "root")
    require(config.get("schema") == "axon.card.v1", "schema must be axon.card.v1")
    meta = config.get("meta")
    keys(meta, META, "meta")
    require(identifier(meta.get("id")), "meta.id must be a protocol ID")
    require(meta["id"] != "synapse_core", "synapse_core is reserved")
    require(type(meta.get("ttl")) is int and meta["ttl"] > 0, "meta.ttl must be a positive integer")
    for key, value in meta.items():
        if key == "ttl": continue
        if key == "tags": require(isinstance(value, list) and all(isinstance(x, str) for x in value), "meta.tags must be strings")
        else: require(isinstance(value, str), f"meta.{key} must be a string")
    layout = config.get("layout")
    keys(layout, {"type", "section"}, "layout")
    require(layout.get("type") == "sections", "layout.type must be sections")
    sections = layout.get("section")
    require(isinstance(sections, list), "layout.section must be an array")
    definitions = config.get("components")
    require(isinstance(definitions, dict), "components must be a table")
    components, actions = {}, set()
    for cid, definition in definitions.items():
        require(identifier(cid), "invalid component ID")
        require(isinstance(definition, dict), "component must be a table")
        kind = definition.get("type")
        require(isinstance(kind, str) and kind in PROPS, "unsupported component type")
        keys(definition, COMMON | PROPS[kind], f"components.{cid}")
        require(isinstance(definition.get("label"), str), "component label is required")
        component = copy.deepcopy(definition)
        component["id"] = cid
        default = component.pop("default", None)
        for field in ("unit", "uri", "text"):
            if field in component: require(isinstance(component[field], str), f"{field} must be a string")
        if "copyable" in component: require(type(component["copyable"]) is bool, "copyable must be boolean")
        monitors = component.get("monitors", [])
        require(isinstance(monitors, list), "monitors must be an array")
        for monitor in monitors:
            keys(monitor, {"condition", "severity", "message"}, "monitor")
            require(all(isinstance(monitor.get(x), str) and monitor[x] for x in ("condition", "severity", "message")), "monitor fields must be nonempty strings")
        component["monitors"] = monitors
        if kind == "gauge":
            require(number(component.get("min")) and number(component.get("max")) and component["max"] > component["min"], "gauge needs finite min < max")
            thresholds = component.get("thresholds", {})
            require(isinstance(thresholds, dict) and all(isinstance(v, str) for v in thresholds.values()), "thresholds must map strings to strings")
            default = component["min"] if default is None else default
        elif kind == "status_indicator":
            mapping = component.get("mapping")
            require(isinstance(mapping, dict) and bool(mapping), "status_indicator needs a mapping")
            for state in mapping.values():
                keys(state, {"text", "color", "icon", "animate"}, "mapping state")
                require(all(isinstance(state.get(x), str) for x in ("text", "color", "icon")), "mapping state requires text/color/icon")
                if "animate" in state: require(type(state["animate"]) is bool, "animate must be boolean")
            require(isinstance(default, str) and default in mapping, "status_indicator default must have a mapping")
        elif kind == "log_stream":
            limit = component.get("max_items", 10)
            require(type(limit) is int and limit > 0, "max_items must be a positive integer")
            component["max_items"] = limit
            default = [] if default is None else default
        elif kind == "action_group":
            require(default is None, "action_group cannot have a default")
            items = component.get("items")
            require(isinstance(items, list), "action_group requires items")
            for item in items:
                keys(item, {"id", "label", "style", "confirm"}, "action item")
                aid = item.pop("id", None)
                require(identifier(aid) and aid not in actions, "action IDs must be valid and unique")
                require(isinstance(item.get("label"), str), "action label is required")
                if "style" in item: require(isinstance(item["style"], str), "action style must be string")
                if "confirm" in item: require(type(item["confirm"]) is bool, "action confirm must be boolean")
                actions.add(aid)
                item["action_id"] = aid
        elif kind == "link":
            require(default is None, "link cannot have a default")
            require(isinstance(component.get("uri"), str) and component["uri"].startswith(("https://", "http://")), "link uri must use http or https")
        if kind not in ("action_group", "link"):
            validate_value(component, default)
            if kind == "log_stream":
                default = [default] if isinstance(default, str) and default else ([] if default == "" else default)
                default = default[-component["max_items"]:]
            component["value"] = default
        components[cid] = component
    roots = []
    for section in sections:
        keys(section, {"title", "components"}, "layout section")
        refs = section.get("components")
        require(isinstance(section.get("title", ""), str), "section title must be string")
        require(isinstance(refs, list) and all(isinstance(cid, str) and cid in components for cid in refs), "layout references unknown component")
        roots.append({"type": "section", "title": section.get("title", ""), "children": refs})
    referenced = {cid for section in roots for cid in section["children"]}
    for cid in components.keys() - referenced:
        warnings.warn(f"Component {cid} is not referenced by layout", stacklevel=2)
    return dict(meta, api_version="v1", status="online", message="", layout={"type": "sections", "root": roots}, components=components), actions
