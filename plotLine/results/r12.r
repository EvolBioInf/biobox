suppressPackageStartupMessages(require(ggplot2)) ||
    stop("no support for ggplot2")
data <- read.table(file="/tmp/tmp_1070323159.dat")
plot <- ggplot(data, aes(V1, V2, group=V3))

plot <- plot + xlab(NULL)


plot <- plot + ylab(NULL)


plot <- plot + scale_y_log10()

plot <- plot + geom_path(aes(color=factor(V3)))

plot <- plot + labs(color="")

x11()
plot(plot)
while(names(dev.cur()) != 'null device')
    Sys.sleep(0.01)
