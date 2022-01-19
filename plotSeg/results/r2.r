library("ggplot2")
d <- read.table(file="stdin")
p <- ggplot(d, aes(V1, V2, xend=V3, yend=V4))
p <- p + geom_segment()
p <- p + xlab("x")

p <- p + ylab(NULL)



x11(
)
plot(p)
while(names(dev.cur()) != 'null device')
    Sys.sleep(0.1)
