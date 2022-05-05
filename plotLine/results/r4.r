library("ggplot2")
data <- read.table(file="stdin")
plot <- ggplot(data, aes(V1, V2, group=V3))

plot <- plot + xlab(NULL)


plot <- plot + ylab(NULL)



plot <- plot + geom_path(aes(color=factor(V3)))

plot <- plot + labs(color="")

plot <- plot + geom_point()
x11()
plot(plot)
while(names(dev.cur()) != 'null device')
    Sys.sleep(1)
