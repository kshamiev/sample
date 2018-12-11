package icons // import "application/controllers/resource/icons"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"gopkg.in/webnice/web.v1/mime"
)

var (
	/*
	 Documentation: https://www.w3.org/2005/10/howto-favicon
	 Documentation: https://stackoverflow.com/questions/2268204/favicon-dimensions
	 Documentation: https://en.wikipedia.org/wiki/Favicon
	 <link rel="icon" type="image/png" href="/favicon.png" />
	 <link rel="shortcut icon" href="/favicon.ico" />
	 <link rel="icon" type="image/vnd.microsoft.icon" href="/favicon.ico" />
	 <link rel="icon" type="image/png" href="/favicon-16x16.png" sizes="16x16" />
	 <link rel="icon" type="image/png" href="/favicon-32x32.png" sizes="32x32" />
	*/
	sizeFavicon = []size{
		{0, 0, "ico", mime.ImageICO, false},
		{0, 0, "png", mime.ImagePNG, false},
		{16, 16, "png", mime.ImagePNG, false},
		{32, 32, "png", mime.ImagePNG, false},
		{36, 36, "png", mime.ImagePNG, false},
		{48, 48, "png", mime.ImagePNG, false},
		{70, 70, "png", mime.ImagePNG, false},
		{72, 72, "png", mime.ImagePNG, false},
		{96, 96, "png", mime.ImagePNG, false},
		{128, 128, "png", mime.ImagePNG, false},
		{144, 144, "png", mime.ImagePNG, false},
		{150, 150, "png", mime.ImagePNG, false},
		{192, 192, "png", mime.ImagePNG, false},
		{196, 196, "png", mime.ImagePNG, false},
		{256, 256, "png", mime.ImagePNG, false},
		{310, 150, "png", mime.ImagePNG, true},
		{310, 310, "png", mime.ImagePNG, false},
		{460, 460, "png", mime.ImagePNG, false},
		{512, 512, "png", mime.ImagePNG, false},
	}

	/*
	 Documentation: https://developer.apple.com/library/archive/documentation/AppleApplications/Reference/SafariWebContent/pinnedTabs/pinnedTabs.html
	 <meta name="apple-mobile-web-app-title" content="WEB DESK" />
	 <meta name="apple-mobile-web-app-status-bar-style" content="black" />
	 <link rel="apple-touch-icon" href="/apple-touch-icon.png" />
	 <link rel="mask-icon" href="/safari-pinned-tab.svg" color="#FFFFFF" />
	 <link rel="apple-touch-startup-image" href="/apple-touch-startup-image.png" />
	 <link rel="apple-touch-icon" sizes="57x57" href="/apple-touch-icon-57x57.png" />
	 <link rel="apple-touch-icon" sizes="60x60" href="/apple-touch-icon-60x60.png" />
	 <link rel="apple-touch-icon-precomposed" sizes="57x57" href="/apple-touch-icon-57x57.png" />
	 <link rel="apple-touch-icon-precomposed" sizes="60x60" href="/apple-touch-icon-60x60.png" />
	*/
	sizeAppleTouchIcon = []size{
		{0, 0, "png", mime.ImagePNG, false},
		{57, 57, "png", mime.ImagePNG, false},
		{60, 60, "png", mime.ImagePNG, false},
		{72, 72, "png", mime.ImagePNG, false},
		{76, 76, "png", mime.ImagePNG, false},
		{114, 114, "png", mime.ImagePNG, false},
		{120, 120, "png", mime.ImagePNG, false},
		{144, 144, "png", mime.ImagePNG, false},
		{152, 152, "png", mime.ImagePNG, false},
		{180, 180, "png", mime.ImagePNG, false},
		{320, 320, "png", mime.ImagePNG, false},
		{320, 460, "png", mime.ImagePNG, true},
		{640, 640, "png", mime.ImagePNG, false},
		{640, 920, "png", mime.ImagePNG, true},
		{640, 1096, "png", mime.ImagePNG, true},
		{748, 1024, "png", mime.ImagePNG, true},
		{768, 1004, "png", mime.ImagePNG, true},
		{1496, 2048, "png", mime.ImagePNG, true},
		{1536, 2008, "png", mime.ImagePNG, true},
		{2048, 2048, "png", mime.ImagePNG, false},
	}

	/*
	 // Documentation: https://developers.google.com/web/fundamentals/web-app-manifest/
	 <link rel="manifest" href="/manifest.json" />
	 <meta name="theme-color" content="#FFFFFF" />
	*/
	sizeAndroidChrome = []size{
		{36, 36, "png", mime.ImagePNG, false},
		{48, 48, "png", mime.ImagePNG, false},
		{72, 72, "png", mime.ImagePNG, false},
		{96, 96, "png", mime.ImagePNG, false},
		{144, 144, "png", mime.ImagePNG, false},
		{192, 192, "png", mime.ImagePNG, false},
		{256, 256, "png", mime.ImagePNG, false},
	}

	/*
	 // Documentation: https://docs.microsoft.com/en-us/previous-versions/windows/internet-explorer/ie-developer/platform-apis/dn320426(v=vs.85)
	 // Documentation: https://docs.microsoft.com/en-us/previous-versions/windows/internet-explorer/ie-developer/samples/dn455106(v%3dvs.85)
	 <meta name="msapplication-TileColor" content="#FFFFFF" />
	 <meta name="application-name" content="WEB DESK">
	 <meta name="msapplication-starturl" content="/?source=pwa">
	 <meta name="msapplication-config" content="/browserconfig.xml" /> или <meta name="msapplication-config" content="none" />
	 <meta name="msapplication-TileImage" content="/mstile-144x144.png" />
	 <meta name="msapplication-square70x70logo" content="/mstile-70x70.png" />
	 <meta name="msapplication-square150x150logo" content="/mstile-150x150.png" />
	 <meta name="msapplication-wide310x150logo" content="/mstile-310x150.png" />
	 <meta name="msapplication-square310x310logo" content="/mstile-310x310.png" />
	*/
	sizeMstile = []size{
		{70, 70, "png", mime.ImagePNG, false},
		{144, 144, "png", mime.ImagePNG, false},
		{150, 150, "png", mime.ImagePNG, false},
		{310, 150, "png", mime.ImagePNG, true},
		{310, 310, "png", mime.ImagePNG, false},
	}
)
