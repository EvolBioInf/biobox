library("ggplot2")
d <- read.table(file="stdin")
p <- ggplot(d, aes(V1, V2, xend=V3, yend=V4))
p <- p + theme_void()
p <- p + geom_segment()
p <- p + geom_text(aes(label=as.character(V5),
                         angle=V6, hjust=V7, vjust=V8))
p <- p + theme(plot.title=element_text(hjust=0.5))
p <- p + ggtitle("newick_1")
p <- p + xlim(-0.1, 1.1)
p <- p + ylim(-0.4, 5.2)

x11()
plot(p)
while(names(dev.cur()) != 'null device')
    Sys.sleep(1)

