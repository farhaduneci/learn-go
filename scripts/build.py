#!/usr/bin/env python3
"""Render content/**/*.md into lessons/*.html, reference/*.html and the manifest.

Each source file starts with a front-matter block:

    ---
    title: Toolchain & modules
    section: Part 1 · Go core
    summary: One line for the learning-path card (lessons only).
    description: The lesson's opening paragraph.
    status: upcoming            # lessons only: done | in-progress | next | upcoming
    source: <a href="...">Name</a> — what to read.
    quiz: Question text? || answer
    ---

Everything after the block is Markdown (converted with pandoc). Lesson order,
titles, prev/next links and the index manifest all come from these files, so
nothing is typed twice.

Usage: python3 scripts/build.py   (requires pandoc)
"""

import html
import json
import pathlib
import subprocess

ROOT = pathlib.Path(__file__).resolve().parent.parent
HEAD = """<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<link rel="preload" href="../assets/fonts/InterVariable.woff2" as="font" type="font/woff2" crossorigin />
<title>{title}</title>
<script>
(function () {{
  try {{
    if (localStorage.getItem("theme") === "light") {{
      document.documentElement.setAttribute("data-theme", "light");
    }}
  }} catch (e) {{}}
}})();
</script>
<link rel="icon" type="image/svg+xml" href="../assets/favicon.svg" />
<link rel="alternate icon" href="../assets/favicon.ico" />
<link rel="stylesheet" href="../assets/style.css" />
</head>
<body>

<button class="theme-toggle" id="theme-toggle" type="button" aria-label="Switch to light theme"><i class="icon" data-icon="moon" aria-hidden="true"></i></button>
"""
TOC = """<nav class="toc" data-toc>
  <div class="toc-title"><i class="icon" data-icon="list-bullets" aria-hidden="true"></i>On this page</div>
</nav>
"""
ASK = """<div class="ask-agent">
  <i class="icon" data-icon="chat-circle-dots" aria-hidden="true"></i>
  <span>Stuck, or want a deeper dive on one idea? Ask a coding agent to explain it with a different example, or to review the exercise solution line by line.</span>
</div>
"""
BACK = '<a href="../index.html"><i class="icon" data-icon="arrow-left" aria-hidden="true"></i>Learning Path</a>'
ARROW = '<i class="icon" data-icon="arrow-right" aria-hidden="true"></i>'


def parse(path):
    text = path.read_text()
    _, front, body = text.split("---\n", 2)
    meta = {"quiz": []}
    for line in front.strip().splitlines():
        key, _, value = line.partition(":")
        key, value = key.strip(), value.strip()
        if key == "quiz":
            q, _, a = value.partition("||")
            meta["quiz"].append({"q": q.strip(), "a": a.strip()})
        else:
            meta[key] = value
    meta["slug"] = path.stem
    meta["body"] = pandoc(body)
    return meta


def pandoc(markdown):
    return subprocess.run(
        ["pandoc", "-f", "markdown-smart", "-t", "html", "--no-highlight", "--wrap=none"],
        input=markdown, capture_output=True, text=True, check=True,
    ).stdout


def quiz_attr(quiz):
    # q is inserted with innerHTML, so HTML-escape it; a is compared as plain text.
    data = [{"q": html.escape(item["q"], quote=False), "a": item["a"]} for item in quiz]
    return html.escape(json.dumps(data, ensure_ascii=False), quote=True).replace("&quot;", '"').replace("'", "&#39;")


def lesson_page(n, meta, prev, nxt):
    parts = [HEAD.format(title=f"Lesson {n}: {html.escape(meta['title'])}")]
    parts.append(f'<p class="lesson-kicker">{BACK} &middot; Lesson {n} &middot; {meta["section"]}</p>')
    parts.append(f"<h1>{html.escape(meta['title'])}</h1>\n")
    parts.append(f"<p>{meta['description']}</p>\n")
    parts.append(TOC)
    parts.append('<div class="callout">\n  <div class="callout-title"><i class="icon" data-icon="book-open-text" aria-hidden="true"></i>Primary source</div>\n  ' + meta["source"] + "\n</div>\n")
    parts.append(meta["body"])
    if meta["quiz"]:
        parts.append('<h2 id="retrieval-check">Retrieval check</h2>\n<p>No scrolling up. Short answers.</p>\n')
        parts.append(f"<div class=\"quiz\" data-questions='{quiz_attr(meta['quiz'])}'></div>\n")
    parts.append(ASK)
    prev_html = f'<a href="{prev[0]}"><i class="icon" data-icon="arrow-left" aria-hidden="true"></i>{prev[1]}</a>' if prev else BACK
    next_html = f'<a href="{nxt[0]}">Next: {nxt[1]}{ARROW}</a>' if nxt else f'<a href="../capstone/README.md">Next: the capstone{ARROW}</a>'
    parts.append(f'<div class="lesson-nav">\n  <span>{prev_html}</span>\n  <span class="next">{next_html}</span>\n</div>\n')
    parts.append(FOOT.format(quiz='<script src="../assets/quiz.js"></script>\n' if meta["quiz"] else ""))
    return "\n".join(parts)


def reference_page(meta):
    parts = [HEAD.format(title=f"Reference: {html.escape(meta['title'])}")]
    parts.append(f'<p class="lesson-kicker">{BACK} &middot; Reference &middot; Quick Lookup</p>')
    parts.append(f"<h1>{html.escape(meta['title'])}</h1>\n")
    parts.append(f'<p class="citation">{meta["source"]}</p>\n')
    parts.append(TOC)
    parts.append(meta["body"])
    parts.append(f'<div class="lesson-nav">\n  <span>{BACK}</span>\n  <span></span>\n</div>\n')
    parts.append(FOOT.format(quiz=""))
    return "\n".join(parts)


FOOT = """<script src="../assets/icons.js"></script>
<script src="../assets/toc.js"></script>
{quiz}<script src="../assets/theme.js"></script>
</body>
</html>
"""


def js(value):
    return json.dumps(value, ensure_ascii=False)


def main():
    lessons = [parse(p) for p in sorted((ROOT / "content/lessons").glob("*.md"))]
    refs = [parse(p) for p in sorted((ROOT / "content/reference").glob("*.md"))]

    for i, meta in enumerate(lessons):
        n = i + 1
        prev = (f"{lessons[i-1]['slug']}.html", f"Lesson {n-1}") if i > 0 else None
        nxt = (f"{lessons[i+1]['slug']}.html", html.escape(lessons[i + 1]["title"])) if i + 1 < len(lessons) else None
        (ROOT / "lessons" / f"{meta['slug']}.html").write_text(lesson_page(n, meta, prev, nxt))
    for meta in refs:
        (ROOT / "reference" / f"{meta['slug']}.html").write_text(reference_page(meta))

    out = ["// Generated by scripts/build.py from content/ front matter. Do not edit by hand.",
           '// status: "done" | "in-progress" | "next" | "upcoming"', "window.LESSONS = ["]
    for i, meta in enumerate(lessons):
        out.append("  {")
        out.append(f"    title: {js('Lesson ' + str(i + 1) + ' &middot; ' + html.escape(meta['title']))},")
        out.append(f"    description: {js(meta['section'] + ' — ' + meta.get('summary', meta['description']))},")
        out.append(f"    href: {js('lessons/' + meta['slug'] + '.html')},")
        out.append(f"    status: {js(meta.get('status', 'upcoming'))},")
        out.append("  },")
    out.append("];\n\nwindow.REFERENCES = [")
    for meta in refs:
        out.append("  {")
        out.append(f"    title: {js('Reference: ' + meta['title'])},")
        out.append(f"    description: {js(meta['description'])},")
        out.append(f"    href: {js('reference/' + meta['slug'] + '.html')},")
        out.append("  },")
    out.append("];\n")
    (ROOT / "assets/lessons-manifest.js").write_text("\n".join(out))
    print(f"built {len(lessons)} lessons, {len(refs)} references")


if __name__ == "__main__":
    main()
