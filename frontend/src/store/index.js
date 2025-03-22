import { defineStore } from 'pinia';

const useStore = defineStore('store', {
    persist: {
        enabled: true,
    },
    state: () => ({
        chats: {
            chatType: "1",
            stream: true,
            jsonFormat: false,
            temperature: 7,
            topK: 0,
            topP: 7.5,
            sourceType: "1",
            cacheLimit: 6,
            cnAnswer: true,
            model: "",
            similarityScore: 7,
            result: [],
            filename: "",
            flowId: "",
            mode: "1",
            pptTemplate: "",
            pageSize: 10,
            uploadFileId: "",
        },
        modeType: "1",
        channel: null,
        isAbout: true,
        loading: true,
        forms: {
            action: "",
        },
        queryData: {},
        planId: {}
    }),
    getters: {

    },
    actions: {
        setAction(value) {
            this.forms.action = value;
        },
        setChannel(value) {
            this.channel = value;
        },
        setModeType(value) {
            this.modeType = value;
        },
        setIsAbout(value) {
            this.isAbout = value;
        },
        setQueryData(key, value) {
            this.queryData[key] = value;
        },
        cleanResult(v) {
            this.chats.result = [];
        },
        setResult(v) {
            let idx = this.chats.result.length - 1;
            if (idx < 0) {
                idx = 0;
            }
            this.chats.result[idx] += v;
        },
        addResult(v) {
            this.chats.result.push(v);
            this.chats.result.push("");
            
        },
        setChats(v) {
            this.chats = v;
        },
    }
})
export default useStore