import os
import re
import json
import shutil
from pathlib import Path
from jinja2 import Template
from titlecase import titlecase
from urllib.parse import quote


def clear_output_dir(output_dir):
    if os.path.exists(output_dir):
        shutil.rmtree(output_dir)
    os.makedirs(output_dir, exist_ok=True)


def ensure_symlink(target, link_name):
    # Remove existing symlink or file; avoid os.remove on dirs.
    if os.path.islink(link_name) or os.path.isfile(link_name):
        os.remove(link_name)
    elif os.path.isdir(link_name):
        shutil.rmtree(link_name)
    # Use relative targets to keep the repo portable.
    rel_target = os.path.relpath(target, start=os.path.dirname(link_name))
    os.symlink(rel_target, link_name)
    print(f"Linked {link_name} → {rel_target}")


def url_encode(value):
    if not isinstance(value, str):
        value = str(value)
    return quote(value, safe="")


def safe_filename(name):
    return titlecase(re.sub(r"[^\w\-_ ]", "", name.strip()))


def find_matching_jobs(json_dir="json", md_dir="md"):
    jobs = []
    json_files = {
        os.path.splitext(f)[0]: os.path.join(json_dir, f)
        for f in os.listdir(json_dir)
        if f.endswith(".json")
    }
    md_files = {
        os.path.splitext(f)[0]: os.path.join(md_dir, f)
        for f in os.listdir(md_dir)
        if f.endswith(".md")
    }
    for table in sorted(json_files):
        if table in md_files:
            jobs.append((json_files[table], md_files[table], table))
    return jobs


def process_json_to_md(json_file, template_file, output_dir, feature_count=7):
    with open(template_file, encoding="utf-8-sig") as f:
        template = Template(f.read())
    with open(json_file, encoding="utf-8-sig") as jf:
        data = json.load(jf)
        for row in data:
            content = template.render(**row, url_encode=url_encode)
            content = content.strip() + "\n"
            filename = safe_filename(row["name"])  # assuming JSON keys are lower
            outpath = f"{output_dir}/{filename}.md"
            with open(outpath, "w", encoding="utf-8-sig") as outfile:
                outfile.write(content)
                print(outpath)


if __name__ == "__main__":
    base_dir = Path(__file__).resolve().parent
    repo_root = base_dir.parent
    docs_dir = base_dir / "docs"

    jobs = find_matching_jobs()
    for json_file, template_file, table in jobs:
        output_dir = repo_root / table
        clear_output_dir(output_dir)
        process_json_to_md(json_file, template_file, output_dir)

        link_name = docs_dir / table
        ensure_symlink(output_dir, link_name)
