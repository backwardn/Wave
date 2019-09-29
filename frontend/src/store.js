import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import router from "./router.js"

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    authenticationState: "unauthenticated",
    loggedIn: true,
    version: 1,
    currentUser: {}
  },
  getters: {
    loggedIn: state => state.loggedIn,
    authenticationState: state => state.authenticationState,
    version: state => state.version,
    currentUser: state => state.currentUser
  },
  actions: {
    authenticate ({commit}, credentials) {
     return axios({url: '/sessions/create', data: credentials, method: 'POST', crossdomain: true, withCredentials: true })
      .then((resp) => {
        commit("setCurrentUser", resp.data)
        commit('authSuccess', credentials.username)
      })
      .catch(err => {
        commit('authFailed', err)
      })
    },
    logout ({commit}) {
     return axios({url: '/sessions/destroy', method: 'POST', crossdomain: true, withCredentials: true })
      .then(() => {
        commit('logout')
      })
    },
    settings ({commit}) {
      commit("settings")
    },
    dashboard ({commit}) {
      commit("dashboard")
    },
    setEnvironment ({commit}) {
      axios({url: '/status', method: 'GET', crossdomain: true, withCredentials: true })
        .then((resp) => {
          commit("setCurrentUser", resp.data)
        })
        .catch(() => {
          commit('logout')
        })
    }
  },
  mutations: {
    setVersion: (state, newVersion) => {
      state.version = newVersion
    },
    authRequest (state) {
      state.authenticationState = "loading"
    },
    setCurrentUser: (state, user) => {
      state.currentUser = user
    },
    authSuccess (state) {
      state.authenticationState = "success"
      state.loggedIn = true
      router.push('/')
    },
    authFailed (state) {
      state.authenticationState = "failed"
    },
    logout (state) {
      state.loggedIn = false
      state.authenticationState = "logged_out"
      router.push("/login")
    },
    settings () {
      router.push("/settings")
    },
    dashboard () {
      router.push("/")
    }
  }
})
