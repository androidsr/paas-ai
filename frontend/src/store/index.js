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
            result: "",
            filename: "",
            flowId: "",
            mode: "1",
            pptTemplate: "",
            pageSize: 10,
            uploadFileId: "",
        },
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
        setIsAbout(value) {
            this.isAbout = value;
        },
        setQueryData(key, value) {
            this.queryData[key] = value;
        },
        setResult(v) {
            this.chats.result = v;
        },
        addResult(v) {
            this.chats.result += v;
        },
        setChats(v) {
            this.chats = v;
        },
    }
})
export default useStore