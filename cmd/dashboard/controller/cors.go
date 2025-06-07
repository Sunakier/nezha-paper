package controller

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Sunakier/nezha-paper/service/singleton"
)

// corsMiddleware handles CORS for API requests
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			// No Origin header, proceed without CORS headers
			c.Next()
			return
		}

		// Always allow OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400") // 24 hours
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Check if it's same origin
		if checkSameOrigin(c.Request) {
			c.Next()
			return
		}

		// Check if we have allowed origins configured for WebSocket
		if singleton.Conf.WSAllowOrigins != "" {
			// Create a map of allowed origins for faster lookup
			allowedOrigins := make(map[string]bool)
			for _, allowedOrigin := range strings.Split(singleton.Conf.WSAllowOrigins, ",") {
				allowedOrigin = strings.TrimSpace(allowedOrigin)
				if allowedOrigin != "" {
					allowedOrigins[allowedOrigin] = true
				}
			}

			// Check if origin is in the allowed list
			if origin != "" {
				u, err := url.Parse(origin)
				if err == nil {
					if allowedOrigins[u.Host] || allowedOrigins[origin] {
						c.Header("Access-Control-Allow-Origin", origin)
						c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
						c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
						c.Header("Access-Control-Allow-Credentials", "true")
						c.Next()
						return
					}
				}
			}
		} else if singleton.Conf.Debug {
			// Allow CORS from loopback addresses in debug mode
			u, err := url.Parse(origin)
			if err == nil {
				host := u.Hostname()
				if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
					c.Header("Access-Control-Allow-Origin", origin)
					c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
					c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
					c.Header("Access-Control-Allow-Credentials", "true")
					c.Next()
					return
				} else {
					// Handle domains like "localhost"
					ip, err := net.LookupHost(host)
					if err == nil && len(ip) > 0 {
						if netIP := net.ParseIP(ip[0]); netIP != nil && netIP.IsLoopback() {
							c.Header("Access-Control-Allow-Origin", origin)
							c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
							c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
							c.Header("Access-Control-Allow-Credentials", "true")
							c.Next()
							return
						}
					}
				}
			}
		}

		// If we get here, CORS is not allowed
		c.Next()
	}
}
