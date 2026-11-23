"""
pii_strategies_demo.py — PII Detection: All 4 Strategies Explained
====================================================================
Runs entirely with the langchain library — no AWS credentials needed.
Uses the internal _redaction helpers directly so you can see exactly
what each strategy does to a piece of text, step by step.

Run:
    python pii_strategies_demo.py

Built-in PII types supported by langchain out of the box:
  email, credit_card, ip, mac_address, url

Custom types are also supported via your own regex (shown below).
"""

import re
from langchain.agents.middleware._redaction import (
    apply_strategy,
    PIIDetectionError,
    detect_email,
    detect_credit_card,
    detect_ip,
    detect_mac_address,
    detect_url,
)

# ── Helpers ──────────────────────────────────────────────────────────────────

DIVIDER = "─" * 62

def section(title: str):
    print(f"\n{DIVIDER}")
    print(f"  {title}")
    print(DIVIDER)

def show(label: str, value: str):
    print(f"  {label:<12} {value}")


# ── Sample texts ──────────────────────────────────────────────────────────────

TEXT_BASIC = (
    "Hi, I'm John. "
    "Email me at john@avasoft.com, "
    "my card is 5105-1051-0510-5100, "
    "and I'm at IP 192.168.0.1."
)

TEXT_APIKEY = "My OpenAI key is sk-abcdefghijklmnopqrstuvwxyz123456 — keep it secret."

TEXT_ALL = (
    "Contact: alice@corp.io | "
    "Card: 4111-1111-1111-1111 | "
    "IP: 10.0.0.5 | "
    "MAC: 00:1A:2B:3C:4D:5E | "
    "Docs: https://internal.corp.io/reports"
)


# ─────────────────────────────────────────────────────────────────────────────
# 1. DETECT — find all PII spans in a string
# ─────────────────────────────────────────────────────────────────────────────
# Each detector returns a list of dicts:
#   { 'type': str, 'value': str, 'start': int, 'end': int }
# 'start' and 'end' are character positions — used by apply_strategy to
# slice and replace the exact substring.

section("STEP 1 — Detect PII spans (no changes yet)")
print(f"\n  Input: {TEXT_BASIC}\n")

matches = (
    detect_email(TEXT_BASIC)
    + detect_credit_card(TEXT_BASIC)
    + detect_ip(TEXT_BASIC)
)

# for m in matches:
#     print(m,"MATCHES")
#     print(f"  [{m['type']:15}]  value={m['value']!r:30}  pos={m['start']}:{m['end']}")


# ─────────────────────────────────────────────────────────────────────────────
# 2. REDACT — replace with a labelled placeholder
# ─────────────────────────────────────────────────────────────────────────────
# Format: [REDACTED_<TYPE>]
# → The original value is completely gone.
# → The model still understands the *structure* of the sentence
#   ("there was an email here") without seeing the real data.
# Best for: audit logs, HIPAA/PCI compliance, any stored record.

section("STRATEGY: redact  →  [REDACTED_TYPE]")
print(f"\n  Before : {TEXT_BASIC}")
redacted = apply_strategy(TEXT_BASIC, matches, "redact")
show("After  :", redacted)
# print("""
# #   How it works:
# #     Each detected span is replaced with [REDACTED_<TYPE>].
# #     The rest of the sentence is unchanged.
# #     The model never sees the real value — it only sees the tag.

# #   Use when:
# #     • You want a permanent, auditable record that PII was present
# #       but must not be stored.
# #     • Logs for healthcare or financial systems (HIPAA, PCI-DSS).
# #     • Customer-service agents where support reps should NOT see
# #       full card numbers.
# # """)


# ─────────────────────────────────────────────────────────────────────────────
# 3. MASK — keep the shape, hide the core
# ─────────────────────────────────────────────────────────────────────────────
# Email  : keeps domain   →  john@****.com
# Card   : keeps last 4   →  ****-****-****-5100
# IP     : keeps last octet →  *.*.*.1
# The recognisable fragment lets a human say "yes, that's the card
# ending in 5100" without the full number being exposed.

section("STRATEGY: mask  →  partial obscuring")
print(f"\n  Before : {TEXT_BASIC}")
masked = apply_strategy(TEXT_BASIC, matches, "mask")
show("After  :", masked)
# print("""
#   How it works:
#     Keeps just enough of the value for human recognition
#     (last 4 card digits, email domain, last IP octet)
#     and replaces the rest with ****  or  *.

#   Use when:
#     • Support dashboards where agents need to confirm "which card"
#       without seeing the full number.
#     • Error messages shown to a user: "We can't process the card
#       ending in ****-5100."
#     • Any UI where partial context helps UX but full data is unsafe.
# """)


# ─────────────────────────────────────────────────────────────────────────────
# 4. HASH — deterministic fingerprint
# ─────────────────────────────────────────────────────────────────────────────
# Format: <type_hash:8-char-hex>
# Same input → always the same hash → you can JOIN records later.
# The hash itself reveals nothing about the original value.

section("STRATEGY: hash  →  deterministic fingerprint")
print(f"\n  Before  : {TEXT_BASIC}")
hashed = apply_strategy(TEXT_BASIC, matches, "hash")
show("After   :", hashed)

# Prove determinism: run it twice on the same value
hashed2 = apply_strategy(TEXT_BASIC, matches, "hash")
print(f"\n  Run again on same input → identical? {hashed == hashed2}")
# print("""
#   How it works:
#     SHA-256 is computed over the raw value and truncated to 8 hex chars.
#     The same PII value always produces the same fingerprint — so you can
#     link records ("all requests from <email_hash:855f96e9>") without
#     storing the email itself.

#   Use when:
#     • Analytics: "how many unique users submitted a form?"
#       (count distinct hashes without storing emails)
#     • Linking logs across services without centralising PII.
#     • A/B testing: ensure the same user always goes to the same bucket.

#   NOT suitable for:
#     • Anything where you ever need to recover the original value.
#     • Short or predictable values — a 4-digit PIN hashes to a small
#       search space that can be brute-forced offline.
# """)


# ─────────────────────────────────────────────────────────────────────────────
# 5. BLOCK — hard stop, raise an exception
# ─────────────────────────────────────────────────────────────────────────────
# Unlike the other three (which sanitise and continue), block raises
# PIIDetectionError immediately. Your application layer must catch it.
# Nothing is passed to the model at all.

section("STRATEGY: block  →  raise PIIDetectionError")
print(f"\n  Input: {TEXT_APIKEY}\n")

api_matches = []
pattern = re.compile(r"sk-[a-zA-Z0-9]{32}")
for m in pattern.finditer(TEXT_APIKEY):
    api_matches.append({"type": "api_key", "value": m.group(), "start": m.start(), "end": m.end()})

try:
    apply_strategy(TEXT_APIKEY, api_matches, "block")
except PIIDetectionError as e:
    print(f"  Caught PIIDetectionError: {e}")
    print("""
  How it works:
    Raises immediately — apply_strategy never returns.
    No sanitised text exists. The agent never runs.
    The caller (your application) decides what to do next.

  What your app should do when catching this:
    1. Log the detection event (WITHOUT the PII value itself).
    2. Return a safe message to the user:
       "We detected a sensitive value in your request.
        Please remove it and try again."
    3. Optionally alert a security team.

  Use when:
    • API keys, passwords, SSNs — values that must NEVER enter
      the LLM pipeline under any circumstances.
    • Compliance rules that require hard rejection (not sanitisation).
    • You want the human to fix their input, not silently strip it.
""")


# # ─────────────────────────────────────────────────────────────────────────────
# # 6. COMBINING DETECTORS — all built-in types at once
# # ─────────────────────────────────────────────────────────────────────────────

section("ALL BUILT-IN TYPES: email + card + ip + mac + url")
print(f"\n  Input: {TEXT_ALL}\n")

all_matches = (
    detect_email(TEXT_ALL)
    + detect_credit_card(TEXT_ALL)
    + detect_ip(TEXT_ALL)
    + detect_mac_address(TEXT_ALL)
    + detect_url(TEXT_ALL)
)

for m in all_matches:
    print(f"  detected  [{m['type']:15}]  {m['value']!r}")

print()
show("redact :", apply_strategy(TEXT_ALL, all_matches, "redact"))
show("mask   :", apply_strategy(TEXT_ALL, all_matches, "mask"))
show("hash   :", apply_strategy(TEXT_ALL, all_matches, "hash"))


# ─────────────────────────────────────────────────────────────────────────────
# 7. QUICK COMPARISON TABLE
# ─────────────────────────────────────────────────────────────────────────────

section("SUMMARY: When to use each strategy")
# print("""
#   Strategy │ Original visible? │ Recoverable? │ Linkable? │ Best for
#   ─────────┼───────────────────┼──────────────┼───────────┼─────────────────────
#   redact   │ No  ([REDACTED])  │ No           │ No        │ Compliance logs, HIPAA
#   mask     │ Partial (****1234)│ No           │ No        │ Support UIs, UX messages
#   hash     │ No  (hex digest)  │ No           │ Yes       │ Analytics, record linking
#   block    │ —  (exception)    │ N/A          │ N/A       │ API keys, passwords, SSNs
# """)