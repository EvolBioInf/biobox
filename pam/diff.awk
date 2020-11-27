{
  n = n + 1
  for(i = 2; i < NF; i++)
      a[n][i-1] = $i
}
END {
  maxD = -1
  for(i = 1; i <= 20; i++) {
    for(j = 1; j <= 20; j++) {
      d = a[i][j] - a[j][i]
      if(d < 0)
	d *= -1
      if(d > maxD) {
	maxD = d
	mi = i
	mj = j
      }
    }
  }
  print mi, mj, a[mi][mj], maxD
}
