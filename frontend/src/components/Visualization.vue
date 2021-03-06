<template>
  <div class="columns">
    <div v-if="displayDetails" id="overlay" class="column is-2">
      <div><button class="delete is-pulled-right" v-on:click="hideDetails()"></button></div>
      <div id="details">
        <hr class="spacer" />
        <b>Selected Node:</b>
        <table class="table is-fullwidth">
          <tbody>
            <tr>
              <td><b>MAC</b></td>
              <td>{{ selectedNodeMAC }}</td>
            </tr>
            <tr>
              <td><b>Vendor</b></td>
              <td>{{ selectedNodeVendor }}</td>
            </tr>
            <tr>
              <td><b>Probed For</b></td>
              <td>{{ selectedNodeProbedFor }}</td>
            </tr>
          </tbody>
        </table>
        <hr class="middle" />
        <b>Selected Association:</b>
      </div>
    </div>
    <div id="visualcontainer" class="column">
      <div id="visualization"></div>
    </div>
  </div>
</template>

<script>
  import Vue from 'vue'
  import ForceGraph3D from '3d-force-graph';

  export default {
    name: 'Visualization',
    data: function() {
      return {
        devicesByMAC: new Map(),
        associationsByKey: new Map(),
        isAssociated: new Map(),
        devices: [],
        associations: [],
        onlyShowAssociated: false,
        graph: ForceGraph3D(),
        displayDetails: false,
        selectedNodeMAC: "",
        selectedNodeVendor: "",
        selectedNodeProbedFor: "",
      }
    },
    methods: {
      hideDetails: function() {
        this.displayDetails = false
        var context = this
	document.getElementById("visualcontainer").setAttribute("class", "column")
        Vue.nextTick(function () {
          context.graph.width(document.getElementById("visualization").offsetWidth)
        })
      },
      showDetails: function(node) {
        this.displayDetails = true
        var context = this
	document.getElementById("visualcontainer").setAttribute("class", "column is-10")
        Vue.nextTick(function () {
          context.graph.width(document.getElementById("visualcontainer").offsetWidth)
          context.selectedNodeMAC = node.MAC
          context.selectedNodeVendor = node.Vendor
          context.selectedNodeProbedFor = node.ProbedFor
        })
      },
      updateDevice: function(device) {
        this.devicesByMAC.set(device.MAC, device)
        this.devices = []
        for (var member of this.devicesByMAC.values()) {
          this.devices.push(member)
        }
        this.graph.graphData({
            links: this.associations,
            nodes: this.devices,
        });
      },
      updateAssociation: function(association) {
	this.isAssociated.set(association.source, true)
	this.isAssociated.set(association.target, true)
        this.associationsByKey.set(association.Key, association)
        this.associations = []
        for (var member of this.associationsByKey.values()) {
          this.associations.push(member)
        }
        this.graph.graphData({
            links: this.associations,
            nodes: this.devices,
        });
      },
      nodeFilter: function(node) {
        if (this.onlyShowAssociated) {
          if (this.isAssociated.get(node.MAC)) {
            return true
          } else {
            return false
          }
        }
        return true
      },
    },
    mounted(){
      var element = document.getElementById("visualization")
      this.graph
        .width(element.offsetWidth)
        .nodeVisibility(this.nodeFilter)
        .onNodeClick(node => { this.showDetails(node) })
        .nodeId("MAC")
        .nodeRelSize(6)
        .nodeOpacity(1)
        .linkOpacity(0.8)
        .linkWidth(3)
        .nodeLabel(node => node.MAC)
        .onNodeHover(node => element.style.cursor = node ? 'pointer' : null);

      this.graph.cameraPosition(0,0)
      this.graph(element).graphData({nodes:this.devices, links: this.associations});

      var ws_protocol = "ws://"
      if (window.location.protocol == "https:") {
        ws_protocol = "wss://"
      }
      var socket = new WebSocket(ws_protocol + window.location.host + '/streams/visualizer')

      var context = this
      socket.onmessage = function (event) {
        var msg = JSON.parse(event.data)
        if (msg.type == "NewDevice") {
          context.updateDevice(msg)
        } else if (msg.type == "NewAssociation") {
          context.updateAssociation(msg)
        }
      }
    }
  }
</script>

<style scoped>
  .columns {
    background-color: #000011;
  }
  h1 {
    text-align: center;
  }
  #overlay {
    color: white;
    background-color: #2C2C2C;
    padding: 12px 0px 0px 12px;
  }
  #visualcontainer {
    padding: 12px 0px 0px 0px;
  }
  table {
    background-color: #2C2C2C;
  }
  td {
    color: white;
  }
  hr.spacer {
    visibility: hidden;
    margin-top: 0px;
    margin-bottom: 10px;
  }
  hr.middle {
    margin-top: 5px;
    margin-bottom: 15px;
    visibility: hidden;
  }
</style>
