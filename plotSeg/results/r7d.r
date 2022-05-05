library("ggplot2")
d <- read.table(file="stdin")
p <- ggplot(d, aes(V1, V2, xend=V3, yend=V4))
p <- p + geom_segment()

p <- p + xlab(NULL)

p <- p + ylab(NULL)

p <- p + xlim(100, 500)


p <- p + ylim(100, 500)


quartz(
)
plot(p)
while(names(dev.cur()) != 'null device')
    Sys.sleep(1)
