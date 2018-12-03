cat(".Rprofile: Setting CRAN repositoryn")
r = getOption("repos") # hard code the UK repo for CRAN
r["CRAN"] = "http://cran.r-project.org"
options(repos = r)
rm(r)

install.packages("ggplot2")
install.packages("httr")
install.packages("jsonlite")
install.packages("grid")
install.packages("reshape2")
install.packages("optparse")
