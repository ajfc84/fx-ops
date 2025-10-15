#!/usr/bin/env python3

import sys


def extract_notes(changelog_path, version):
    target_version = version.split('-')[0].split('+')[0]

    with open(changelog_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()

    in_version_section = False
    in_notes_block = False
    notes = []

    version_header = f"### **Version {target_version}**"

    for line in lines:
        line = line.rstrip('\n')

        if line.strip() == version_header:
            in_version_section = True
            continue

        if in_version_section:
            if line.strip() == '---' and not in_notes_block:
                # Início do bloco de notas
                in_notes_block = True
                continue
            elif line.strip() == '---' and in_notes_block:
                # Fim do bloco de notas
                break
            elif in_notes_block:
                notes.append(line)

    if notes:
        print('\n'.join(notes))
    else:
        print(f"Notas para a versão {version} não encontradas.", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Use: extract_notes.py <changelog_path> <version>", file=sys.stderr)
        sys.exit(1)

    changelog_path = sys.argv[1]
    version = sys.argv[2]

    extract_notes(changelog_path, version)
