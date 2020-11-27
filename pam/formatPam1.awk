NR == 1 {
  printf "\t%s", $1
  for(i = 2; i <= NF; i++)
    printf "\t%s", $i
  printf "\n"
}
NR > 1 {
  printf "%s", $1
  for(i = 2; i <= NF; i++)
    printf "\t%.0f", $i*10000
  printf "\n"
}
