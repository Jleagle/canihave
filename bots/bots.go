package bots

func IsBot(ua string) bool {
	var m = map[string]bool{
		"Baiduspider+(+http://www.baidu.com/search/spider.htm":                                                                                                                                                                                               true,
		"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)":                                                                                                                                                                true,
		"Moreoverbot/5.1 (+http://w.moreover.com; webmaster@moreover.com) Mozilla/5.0":                                                                                                                                                                       true,
		"UnwindFetchor/1.0 (+http://www.gnip.com/)":                                                                                                                                                                                                          true,
		"Voyager/1.0":                                                                                                                                                                                                                                        true,
		"PostRank/2.0 (postrank.com)":                                                                                                                                                                                                                        true,
		"R6_FeedFetcher(www.radian6.com/crawler)":                                                                                                                                                                                                            true,
		"R6_CommentReader(www.radian6.com/crawler)":                                                                                                                                                                                                          true,
		"radian6_default_(www.radian6.com/crawler)":                                                                                                                                                                                                          true,
		"Mozilla/5.0 (compatible; Ezooms/1.0; ezooms.bot@gmail.com)":                                                                                                                                                                                         true,
		"ia_archiver (+http://www.alexa.com/site/help/webmasters; crawler@alexa.com)":                                                                                                                                                                        true,
		"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)":                                                                                                                                                                true,
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)":                                                                                                                                                                           true,
		"Mozilla/5.0 (en-us) AppleWebKit/525.13 (KHTML, like Gecko; Google Web Preview) Version/3.1 Safari/525.13":                                                                                                                                           true,
		"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)":                                                                                                                                                                                   true,
		"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)":                                                                                                                                                                            true,
		"Twitterbot/0.1":                                                                                                                                                                                                                                     true,
		"LinkedInBot/1.0 (compatible; Mozilla/5.0; Jakarta Commons-HttpClient/3.1 +http://www.linkedin.com)":                                                                                                                                                 true,
		"bitlybot":                                                                                                                                                                                                                                           true,
		"MetaURI API/2.0 +metauri.com":                                                                                                                                                                                                                       true,
		"Mozilla/5.0 (compatible; Birubot/1.0) Gecko/2009032608 Firefox/3.0.8":                                                                                                                                                                               true,
		"Mozilla/5.0 (compatible; PrintfulBot/1.0; +http://printful.com/bot.html)":                                                                                                                                                                           true,
		"Mozilla/5.0 (compatible; PaperLiBot/2.1)":                                                                                                                                                                                                           true,
		"Summify (Summify/1.0.1; +http://summify.com)":                                                                                                                                                                                                       true,
		"Mozilla/5.0 (compatible; TweetedTimes Bot/1.0; +http://tweetedtimes.com)":                                                                                                                                                                           true,
		"PycURL/7.18.2":                                                                                                                                                                                                                                      true,
		"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)":                                                                                                                                                                          true,
		"Python-urllib/2.6":                                                                                                                                                                                                                                  true,
		"Python-httplib2/$Rev$":                                                                                                                                                                                                                              true,
		"AppEngine-Google; (+http://code.google.com/appengine; appid: lookingglass-server)":                                                                                                                                                                  true,
		"Wget/1.9+cvs-stable (Red Hat modified)":                                                                                                                                                                                                             true,
		"Mozilla/5.0 (compatible; redditbot/1.0; +http://www.reddit.com/feedback)":                                                                                                                                                                           true,
		"Mozilla/5.0 (compatible; MSIE 6.0b; Windows NT 5.0) Gecko/2009011913 Firefox/3.0.6 TweetmemeBot":                                                                                                                                                    true,
		"Mozilla/5.0 (compatible; discobot/1.1; +http://discoveryengine.com/discobot.html)":                                                                                                                                                                  true,
		"Mozilla/5.0 (compatible; Exabot/3.0; +http://www.exabot.com/go/robot)":                                                                                                                                                                              true,
		"Mozilla/5.0 (compatible; SiteBot/0.1; +http://www.sitebot.org/robot/)":                                                                                                                                                                              true,
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1 + FairShare-http://fairshare.cc)":                                                                                                                                                            true,
		"HTTP_Request2/2.0.0beta3 (http://pear.php.net/package/http_request2) PHP/5.3.2":                                                                                                                                                                     true,
		"Mozilla/5.0 (compatible; Embedly/0.2; +http://support.embed.ly/)":                                                                                                                                                                                   true,
		"magpie-crawler/1.1 (U; Linux amd64; en-GB; +http://www.brandwatch.net)":                                                                                                                                                                             true,
		"(TalkTalk Virus Alerts Scanning Engine)":                                                                                                                                                                                                            true,
		"Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)":                                                                                                                                                                            true,
		"Googlebot/2.1 (http://www.googlebot.com/bot.html)":                                                                                                                                                                                                  true,
		"msnbot-NewsBlogs/2.0b (+http://search.msn.com/msnbot.htm)":                                                                                                                                                                                          true,
		"msnbot/2.0b (+http://search.msn.com/msnbot.htm)":                                                                                                                                                                                                    true,
		"msnbot-media/1.1 (+http://search.msn.com/msnbot.htm)":                                                                                                                                                                                               true,
		"Mozilla/5.0 (compatible; oBot/2.3.1; +http://www-935.ibm.com/services/us/index.wss/detail/iss/a1029077?cntxt=a1027244)":                                                                                                                             true,
		"Sosospider+(+http://help.soso.com/webspider.htm)":                                                                                                                                                                                                   true,
		"COMODOspider/Nutch-1.0":                                                                                                                                                                                                                             true,
		"trunk.ly spider contact@trunk.ly":                                                                                                                                                                                                                   true,
		"Mozilla/5.0 (compatible; Purebot/1.1; +http://www.puritysearch.net/)":                                                                                                                                                                               true,
		"Mozilla/5.0 (compatible; MJ12bot/v1.4.0; http://www.majestic12.co.uk/bot.php?+)":                                                                                                                                                                    true,
		"knowaboutBot 0.01":                                                                                                                                                                                                                                  true,
		"Showyoubot (http://showyou.com/support)":                                                                                                                                                                                                            true,
		"Flamingo_SearchEngine (+http://www.flamingosearch.com/bot)":                                                                                                                                                                                         true,
		"MLBot (www.metadatalabs.com/mlbot)":                                                                                                                                                                                                                 true,
		"my-robot/0.1":                                                                                                                                                                                                                                       true,
		"Mozilla/5.0 (compatible; woriobot support [at] worio [dot] com +http://worio.com)":                                                                                                                                                                  true,
		"Mozilla/5.0 (compatible; YoudaoBot/1.0; http://www.youdao.com/help/webmaster/spider/; )":                                                                                                                                                            true,
		"chilitweets.com":                                                                                                                                                                                                                                    true,
		"Mozilla/5.0 (TweetBeagle; http://app.tweetbeagle.com/)":                                                                                                                                                                                             true,
		"OctoBot/2.1 (OctoBot/2.1.0; +http://www.octofinder.com/octobot.html?2.1)":                                                                                                                                                                           true,
		"Mozilla/5.0 (compatible; FriendFeedBot/0.1; +Http://friendfeed.com/about/bot)":                                                                                                                                                                      true,
		"Mozilla/5.0 (compatible; WASALive Bot ; http://blog.wasalive.com/wasalive-bots/)":                                                                                                                                                                   true,
		"Mozilla/5.0 (compatible; Apercite; +http://www.apercite.fr/robot/index.html)":                                                                                                                                                                       true,
		"urlfan-bot/1.0; +http://www.urlfan.com/site/bot/350.html":                                                                                                                                                                                           true,
		"SeznamBot/3.0 (+http://fulltext.sblog.cz/)":                                                                                                                                                                                                         true,
		"Yeti/1.0 (NHN Corp.; http://help.naver.com/robots/)":                                                                                                                                                                                                true,
		"Mozilla/5.0 (Windows; U; Windows NT 6.0; en-GB; rv:1.0; trendictionbot0.4.2; trendiction media ssppiiddeerr; http://www.trendiction.com/bot/; please let us know of any problems; ssppiiddeerr at trendiction.com) Gecko/20071127 Firefox/2.0.0.11": true,
		"yacybot (freeworld/global; amd64 Linux 2.6.35-24-generic; java 1.6.0_20; Asia/en) http://yacy.net/bot.html":                                                                                                                                         true,
		"Mozilla/5.0 (compatible; suggybot v0.01a, http://blog.suggy.com/was-ist-suggy/suggy-webcrawler/)":                                                                                                                                                   true,
		"ssearch_bot (sSearch Crawler; http://www.semantissimo.de)":                                                                                                                                                                                          true,
		"Mozilla/5.0 (compatible; Linux; Socialradarbot/2.0; en-US; crawler@infegy.com)":                                                                                                                                                                     true,
		"wikiwix-bot-3.0":                                                                                                                                                                                                                                    true,
		"Mozilla/5.0 (compatible; AhrefsBot/1.0; +http://ahrefs.com/robot/)":                                                                                                                                                                                 true,
		"Mozilla/5.0 (compatible; DotBot/1.1; http://www.dotnetdotcom.org/, crawler@dotnetdotcom.org)":                                                                                                                                                       true,
		"GarlikCrawler/1.1 (http://garlik.com/, crawler@garik.com)":                                                                                                                                                                                          true,
		"Mozilla/5.0 (compatible; SISTRIX Crawler; http://crawler.sistrix.net/)":                                                                                                                                                                             true,
		"Mozilla/5.0 (compatible; 008/0.83; http://www.80legs.com/webcrawler.html) Gecko/2008032620":                                                                                                                                                         true,
		"PostPost/1.0 (+http://postpo.st/crawlers)":                                                                                                                                                                                                          true,
		"Aghaven/Nutch-1.2 (www.aghaven.com)":                                                                                                                                                                                                                true,
		"SBIder/Nutch-1.0-dev (http://www.sitesell.com/sbider.html)":                                                                                                                                                                                         true,
		"Mozilla/5.0 (compatible; ScoutJet; +http://www.scoutjet.com/)":                                                                                                                                                                                      true,
		"Soup/2011-05-11Z11-51-38–soup–production-2-g251c1f9d/251c1f9d6cdff8491e0b49f4ba3288ec7f3de903 (http://soup.io/)":                                                                                                                                    true,
		"Trapit/1.1":                                                                                                                                                                                                                                         true,
		"Jakarta Commons-HttpClient/3.1":                                                                                                                                                                                                                     true,
		"Readability/0.1":                                                                                                                                                                                                                                    true,
		"kame-rt (support@backtype.com)":                                                                                                                                                                                                                     true,
		"Mozilla/5.0 (compatible; Topix.net; http://www.topix.net/topix/newsfeeds)":                                                                                                                                                                          true,
		"Megite2.0 (http://www.megite.com)":                                                                                                                                                                                                                  true,
		"SkyGrid/1.0 (+http://skygrid.com/partners)":                                                                                                                                                                                                         true,
		"Netvibes (http://www.netvibes.com)":                                                                                                                                                                                                                 true,
		"Zemanta Aggregator/0.7 +http://www.zemanta.com":                                                                                                                                                                                                     true,
		"Owlin.com/1.3 (http://owlin.com/)":                                                                                                                                                                                                                  true,
		"Mozilla/5.0 (compatible; Twitturls; +http://twitturls.com)":                                                                                                                                                                                         true,
		"Tumblr/1.0 RSS syndication (+http://www.tumblr.com/) (support@tumblr.com)":                                                                                                                                                                          true,
		"Mozilla/4.0 (compatible; www.euro-directory.com; urlchecker1.0)":                                                                                                                                                                                    true,
		"Covario-IDS/1.0 (Covario; http://www.covario.com/ids; support at covario dot com)":                                                                                                                                                                  true,
	}

	if _, ok := m[ua]; ok {
		return true
	}
	return false
}
