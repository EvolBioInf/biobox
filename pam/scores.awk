NR == 1 {
  for(i = 1; i <= NF; i++)
    header[i] = $i
}
NR > 1 {
  for(i = 2; i <= NF; i++)
    printf "%c\t%c\t%d\n", $1, header[i-1], $i
}
