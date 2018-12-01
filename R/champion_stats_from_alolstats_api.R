library(httr)
library(jsonlite)
library(grid)
library(reshape2)
library(ggplot2)

options(width=200)

championids <- c()
gameversion = "8.23"

base <- "http://localhost:8000/"

endpoint <- "v1/champions"
call1 <- paste(base,endpoint, sep="")

champions_json <- fromJSON(content(GET(call1), "text"), simplifyMatrix = FALSE, flatten = FALSE)
# print(champions_json)
for (champion in champions_json) {
    for (content in champion) {
        championids <- c(championids, content$key)
    }
}

endpoint <- "v1/stats/champion/byid"
for (championid in championids) {
    call1 <- paste(base,endpoint,"?","id","=", championid, "&gameversion=", gameversion, sep="")
    get_champion_stats_json <- fromJSON(content(GET(call1), "text"), flatten = TRUE)

    get_champion_base.data <- as.data.frame(get_champion_stats_json)
    get_champion_stats.data <- as.data.frame(get_champion_stats_json$lanerolepercentage)

    champName = get_champion_base.data[1,]$championname
    champrealid = get_champion_base.data[1,]$championrealid

    timestamp = get_champion_base.data[1,]$timestamp

    totalCnt = get_champion_base.data[1,]$samplesize
    get_champion_stats.data$winpercentage = (get_champion_stats.data$wins)/get_champion_stats.data$ngames*100.0

    p <- ggplot(data=get_champion_stats.data, aes(x=lane, y=percentage, fill=role)) +
    geom_bar(stat="identity") +
    theme_minimal() + 
        labs(title=paste("Champion Role Distribution:", champName)) +
        labs(y="Percentage [%]", x="Lane") +
    theme(plot.title = element_text(hjust = 0.5)) + 
    theme(text = element_text(size=20)) +
        theme(plot.title=element_text(family="Helvetica", face="bold", size=24)) +
    scale_y_continuous( breaks = scales::pretty_breaks(n = 10), limits=c(0,100) ) + 
    theme(legend.title=element_blank()) +
    theme(        panel.background = element_blank(),
            panel.grid.major = element_blank(), 
            panel.grid.minor = element_blank(),
            axis.line = element_line(colour = "black"),
            panel.border = element_rect(colour = "black", fill=NA, size=1)) +
            geom_text(aes(label=sprintf("WP: %.02f %% (%d)", winpercentage, ngames)), position = position_stack(vjust = 0.5), size = 3) +
            annotation_custom(grid::textGrob(sprintf("Only unranked and ranked PvP games on Summoners Rift with game version %s are considered\nNumber on bars are Win Percentage and games analyzed for that particular role\n%d total games analyzed", gameversion, totalCnt), gp=gpar(col="black", fontsize=10, fontface="italic")), 
                        xmin = -Inf, xmax = Inf, ymin = 100, ymax = 100) +
            annotation_custom(grid::textGrob(sprintf("(%s)", timestamp), gp=gpar(col="darkgrey", fontsize=8, fontface="italic")), 
                        xmin = 0, xmax = 2.3, ymin = -7.5, ymax = 0) 

    cairo_pdf(paste('champion_role_',champrealid,'_',championid,'.pdf', sep=""), width = 10, height = 10/1.2)
        print(p)
    dev.off()

    png(paste('champion_role_',champrealid,'_',championid,'.png', sep=""), width = 10, height = 10/1.2, units = "in", res=300)
        print(p)
    dev.off()

}
