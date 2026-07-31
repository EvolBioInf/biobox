git describe |
    tr '-' ' ' |
    awk '{printf "%s", $1; if($2)printf "-%s", $2; printf "\n"}'
