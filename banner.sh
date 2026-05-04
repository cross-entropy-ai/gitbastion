#!/bin/sh
cat <<'EOF'

    ╔══════════════════════════════════════════════════════════════╗
    ║                                                              ║
    ║    ┌─┐┬┌┬┐┌┐ ┌─┐┌─┐┌┬┐┬┌─┐┌┐┌                             ║
    ║    │ ┬│ │ ├┴┐├─┤└─┐ │ ││ ││││                             ║
    ║    └─┘┴ ┴ └─┘┴ ┴└─┘ ┴ ┴└─┘┘└┘                             ║
    ║                                                              ║
    ║    You have been authenticated successfully.                  ║
    ║    However, interactive shell access is not available.        ║
    ║                                                              ║
    ║    ── Access cluster resources via ──                         ║
    ║                                                              ║
    ║    ProxyJump:                                                 ║
    ║      ssh -J git@<bastion> user@<internal-host>               ║
    ║                                                              ║
    ║    Local Forwarding:                                          ║
    ║      ssh -L 8080:<internal-host>:80 -N git@<bastion>         ║
    ║                                                              ║
    ╚══════════════════════════════════════════════════════════════╝

EOF
exit 0
