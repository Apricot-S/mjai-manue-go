#!/usr/bin/env python3

"""Copy local viewer images and rewrite generated mjai HTML to use them."""

import argparse
import re
import shutil
import sys
from pathlib import Path

SCRIPT_DIR = Path(__file__).resolve().parent
DEFAULT_IMAGE_DIR = SCRIPT_DIR / "images"

TILE_NAMES = (
    [f"ms{i}" for i in range(1, 10)]
    + [f"ps{i}" for i in range(1, 10)]
    + [f"ss{i}" for i in range(1, 10)]
    + ["ji_e", "ji_s", "ji_w", "ji_n", "no", "ji_h", "ji_c", "bk"]
)
RED_TILE_NAMES = ("ms5r", "ps5r", "ss5r")
POSES = (1, 3)


def required_asset_names() -> list[str]:
    names = ["blank.png"]
    for tile in TILE_NAMES:
        for pose in POSES:
            names.append(f"p_{tile}_{pose}.gif")
    for tile in RED_TILE_NAMES:
        for pose in POSES:
            names.append(f"p_{tile}_{pose}.png")
    return sorted(names)


def resource_dir_for(html_path: Path) -> Path:
    text = html_path.read_text(encoding="utf-8")
    match = re.search(r'resourceDir\s*=\s*"([^"]+)"', text)
    if match:
        return html_path.parent / match.group(1)

    candidates = sorted(html_path.parent.glob(f"{html_path.name}.files"))
    if len(candidates) == 1:
        return candidates[0]

    raise ValueError(f"could not determine resource directory for {html_path}")


def copy_images(image_dir: Path, dest_dir: Path) -> None:
    missing = [
        name
        for name in required_asset_names()
        if not (image_dir / name).is_file()
    ]

    if missing:
        details = "\n".join(f"  {name}" for name in missing)
        raise FileNotFoundError(
            f"missing required images in {image_dir}:\n{details}"
        )

    dest_dir.mkdir(parents=True, exist_ok=True)
    for name in required_asset_names():
        shutil.copy2(image_dir / name, dest_dir / name)


def rewrite_archive_player(js_path: Path) -> None:
    text = js_path.read_text(encoding="utf-8")
    old_tile = (
        '"http://gimite.net/mjai/images/p_" + name + "_" + pose + "." + ext'
    )
    new_tile = 'resourceDir + "/images/p_" + name + "_" + pose + "." + ext'
    old_blank = '"http://gimite.net/mjai/images/blank.png"'
    new_blank = 'resourceDir + "/images/blank.png"'

    text = text.replace(old_tile, new_tile)
    text = text.replace(old_blank, new_blank)
    js_path.write_text(text, encoding="utf-8", newline="\n")


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Copy self-match viewer images and rewrite generated mjai HTML resources.",
    )
    parser.add_argument("html", type=Path, help="generated .html file")
    parser.add_argument(
        "--image-dir",
        type=Path,
        default=DEFAULT_IMAGE_DIR,
        help=f"source image directory (default: {DEFAULT_IMAGE_DIR})",
    )
    args = parser.parse_args()

    html_path = args.html.resolve()
    image_dir = args.image_dir.resolve()
    if not html_path.is_file():
        parser.error(f"HTML file does not exist: {html_path}")
    if not image_dir.is_dir():
        parser.error(f"image directory does not exist: {image_dir}")

    try:
        resource_dir = resource_dir_for(html_path)
        js_path = resource_dir / "js" / "archive_player.js"
        if not js_path.is_file():
            raise FileNotFoundError(f"archive_player.js not found: {js_path}")
        copy_images(image_dir, resource_dir / "images")
        rewrite_archive_player(js_path)
    except Exception as exc:
        print(f"replace_assets.py: {exc}", file=sys.stderr)
        return 1

    print(f"rewrote {js_path}")
    print(f"copied images to {resource_dir / 'images'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
