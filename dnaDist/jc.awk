BEGIN {
  m = 0
  inc = 0.01
  while(m < 0.75) {
    j = -3/4*log(1-4/3*m)
    printf "%f\t%f\n", m, j
    m += inc
  }
}
