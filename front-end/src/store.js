import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        config: undefined,
        notifyInfo: undefined,
        notifyError: undefined,
        notifySuccess:undefined,
        account: undefined,
        socketConstants: undefined,
    },
    // 修改全局变量必须通过mutations中的方法
    // mutations只能采用同步方法
    mutations: {
        setAccount(state, payload) {
            state.account = payload;
        },
        setConfig(state, payload) {
            state.config = payload
        },
        setNotifyInfo(state,payload) {
            state.notifyInfo = payload
        },
        setNotifyError(state,payload){
            state.notifyError = payload
        },
        setNotifySuccess(state,payload){
            state.notifySuccess = payload
        },
        setSocketConstants(state,payload) {
            state.socketConstants = payload
        }
    },
    // 异步方法用actions
    // actions不能直接修改全局变量，需要调用commit方法来触发mutation中的方法
    actions: {
    }
});
