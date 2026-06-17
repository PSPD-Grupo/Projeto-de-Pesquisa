from pathlib import Path

from grpc_tools import protoc


ROOT = Path(__file__).resolve().parents[1]
PROJECT_ROOT = ROOT.parents[2]
OUT_DIR = ROOT / "app" / "grpc_generated"

PROTO_DIRS = [
    PROJECT_ROOT / "docs" / "microsservicos" / "ServidorA" / "proto",
    PROJECT_ROOT
    / "docs"
    / "microsservicos"
    / "servidor-b"
    / "src"
    / "main"
    / "proto",
]


def patch_relative_imports(path: Path) -> None:
    content = path.read_text(encoding="utf-8")
    for proto_file in OUT_DIR.glob("*_pb2.py"):
        module = proto_file.stem
        content = content.replace(f"import {module} as ", f"from . import {module} as ")
    path.write_text(content, encoding="utf-8")


def main() -> int:
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    proto_files = [str(path) for proto_dir in PROTO_DIRS for path in proto_dir.glob("*.proto")]

    args = [
        "grpc_tools.protoc",
        *(f"-I{proto_dir}" for proto_dir in PROTO_DIRS),
        f"--python_out={OUT_DIR}",
        f"--grpc_python_out={OUT_DIR}",
        *proto_files,
    ]

    result = protoc.main(args)
    if result != 0:
        return result

    for generated_grpc_file in OUT_DIR.glob("*_pb2_grpc.py"):
        patch_relative_imports(generated_grpc_file)

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
