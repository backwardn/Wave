package controllers

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hkparker/Wave/engines/ids"
	"github.com/hkparker/Wave/engines/visualizer"
	"github.com/hkparker/Wave/models"
)

func createCollector(c *gin.Context) {
	user_info, err := requestJSON(c)
	if err != nil {
		return
	}
	admin, err := currentUser(c)
	if err != nil {
		return
	}

	// Ensure the JSON contains a name
	name, ok := user_info["name"]
	if !ok || name == "" {
		name_error := "no name provided"
		c.JSON(400, gin.H{"error": name_error})
		log.WithFields(log.Fields{
			"at":    "controllers.createCollector",
			"error": name_error,
			"admin": admin.Username,
		}).Error("error creating collector")
		return
	}

	collector, err := models.CreateCollector(name)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.WithFields(log.Fields{
			"at":    "controllers.createCollector",
			"error": err.Error(),
			"admin": admin.Username,
		}).Error("error creating collector")
		return

	}

	c.JSON(200, gin.H{
		"success":     "true",
		"name":        collector.Name,
		"certificate": collector.CaCert,
		"private_key": collector.PrivateKey,
	})
	log.WithFields(log.Fields{
		"at":    "controllers.createCollector",
		"name":  name,
		"admin": admin.Username,
	}).Info("created collector")
}

func getCollectors(c *gin.Context) {
	admin, err := currentUser(c)
	if err != nil {
		return
	}

	collectors, err := models.Collectors()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.WithFields(log.Fields{
			"at":    "controllers.getCollectors",
			"error": err.Error(),
			"admin": admin.Username,
		}).Error("error looking up collectors")
		return
	}
	names := []string{}
	for _, collector := range collectors {
		names = append(names, collector.Name)
	}
	c.JSON(200, names)
}

func deleteCollector(c *gin.Context) {
	user_info, err := requestJSON(c)
	if err != nil {
		return
	}
	admin, err := currentUser(c)
	if err != nil {
		return
	}

	// Ensure the JSON contains a name key
	name, ok := user_info["name"]
	if !ok || name == "" {
		name_error := "no name provided"
		c.JSON(400, gin.H{"error": name_error})
		log.WithFields(log.Fields{
			"at":    "controllers.deleteCollector",
			"error": name_error,
			"admin": admin.Username,
		}).Error("error deleting collector")
		return
	}

	// Ensure the collector exists
	collector, err := models.CollectorByName(name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.WithFields(log.Fields{
			"at":    "controllers.deleteCollector",
			"name":  name,
			"error": err.Error(),
			"admin": admin.Username,
		}).Error("error looking up collector to delete")
		return
	}

	// Delete the collector
	err = collector.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.WithFields(log.Fields{
			"at":    "controllers.deleteCollector",
			"name":  name,
			"error": err.Error(),
			"admin": admin.Username,
		}).Error("error deleting collector")
	} else {
		c.JSON(200, gin.H{"success": "true"})
		log.WithFields(log.Fields{
			"at":    "controllers.deleteCollector",
			"name":  name,
			"admin": admin.Username,
		}).Info("deleted collector")
	}
}

func pollCollector(c *gin.Context) {
	var upgrayedd websocket.Upgrader
	conn, err := upgrayedd.Upgrade(c.Writer, c.Request, nil)
	if err == nil {
		defer conn.Close()
		// Generate an ID for this collector from the certificate presented
		sig := conn.UnderlyingConn().(*tls.Conn).ConnectionState().PeerCertificates[0].Signature
		sha := sha256.Sum256(sig)
		collector_id := hex.EncodeToString(sha[:])

		for {
			_, frame_bytes, err := conn.ReadMessage()
			if err != nil {
				break
			}

			var frame models.Wireless80211Frame
			json.Unmarshal(frame_bytes, &frame)

			visualizer.Insert(frame)
			ids.Insert(string(frame_bytes), frame, collector_id)
		}
	} else {
		log.WithFields(log.Fields{
			"at":    "controllers.pollCollector",
			"error": err.Error(),
		}).Warn("failed to upgrade websocket connection")
	}
}
