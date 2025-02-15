<template>
    <div>
        <a-page-header style="border: 1px solid rgb(235, 237, 240); padding: 8px; margin: 0;" title="函数实现"
            @back="back" />
        <a-card>
            <a-form :model="model" :rules="vRules" :label-col="{ style: { width: '90px' } }">
                <a-form-item label="实现标题" name="title">
                    <a-input v-model:value="model.title" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="函数名称" name="name">
                    <a-input v-model:value="model.name" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="请求地址" name="url">
                    <a-input v-model:value="model.url" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="请求方式" name="method">
                    <a-select v-model:value="model.method" :options="methods" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="数据类型" name="contentType">
                    <a-select v-model:value="model.contentType" :options="contentTypes" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="请求头" name="headers">
                    <a-textarea v-model:value="model.headers" rows="4" :readonly="isReadOnly" />
                </a-form-item>
            </a-form>

            <template #actions v-if="!isReadOnly">
                <a-space :size="16">
                    <a-button @click="back">返回</a-button>
                    <a-button type="primary" @click="submitForm">确认</a-button>
                </a-space>
            </template>
        </a-card>
    </div>
</template>

<script>
import { Add, Edit, Get } from '../../../wailsjs/go/biz/FunctionImplBiz.js';

export default {
    mounted() {
        this.id = this.$route.query.id;
        if (this.id) {
            Get(this.id).then(res => {
                if (res.code == 200) {
                    this.model = res.data || {};
                }
            })
        }
    },
    data() {
        return {
            contentTypes: [
                { label: "application/json", value: "application/json" },
                { label: "application/x-www-form-urlencoded", value: "application/x-www-form-urlencoded" },
                { label: "multipart/form-data", value: "multipart/form-data" },
                { label: "text/html", value: "text/html" },
                { label: "text/xml", value: "text/xml" },
            ],
            methods: [
                { label: "GET", value: "GET" },
                { label: "POST", value: "POST" },
                { label: "PUT", value: "PUT" },
                { label: "DELETE", value: "DELETE" },
            ],
            id: "",
            model: {},
            vRules: {
                funcName: [{ required: true, message: "请输入提示词名称" }],
                funcContent: [{ required: true, message: "请输入提示词内容" }],
                funcContent: [{ required: true, message: "请输入提示词内容" }],
            },
        }
    },
    computed: {
        isReadOnly() {
            return this.$store.forms.action === 'detail';
        }
    },
    methods: {
        submitForm() {
            if (this.$store.forms.action == "add") {
                this.model.id = "";
                Add(this.model).then(res => {
                    if (res.code == 200) {
                        this.$message.success(res.msg);
                        this.back();
                    }
                })
            } else {
                Edit(this.model).then(res => {
                    if (res.code == 200) {
                        this.$message.success(res.msg);
                        this.back();
                    }
                })
            }
        },
        back() {
            this.$router.back();
        }
    }
}
</script>

<style scoped>
.table-header {
    font-weight: bold;
}
</style>
