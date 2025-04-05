<template>
    <div>
        <a-page-header style="border: 1px solid rgb(235, 237, 240); padding: 8px; margin: 0;" title="提示词表单"
            @back="back" />
        <a-card>
            <a-form :model="model" :rules="vRules" :label-col="{ style: { width: '90px' } }">
                <a-form-item label="渠道名称" name="name">
                    <a-input v-model:value="model.name" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="接口地址" name="url">
                    <a-input v-model:value="model.url" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="网页地址" name="originalUrl">
                    <a-input v-model:value="model.originalUrl" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="认证Token" name="token">
                    <a-input v-model:value="model.token" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="支持模型" name="models">
                    <a-input v-model:value="model.models" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="优先级" name="priority">
                    <a-input-number v-model:value="model.priority" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="消息长度" name="maxToken">
                    <a-input-number v-model:value="model.maxToken" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="备注" name="remark">
                    <a-textarea v-model:value="model.remark" :rows="4" :readonly="isReadOnly" />
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
import { Add, Edit, Get } from '../../../wailsjs/go/biz/AiChannelBiz.js';

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
            id: "",
            model: {
                maxToken: -1,
            },
            vRules: {
                name: [{ required: true, message: "请输入渠道名称" }],
                url: [{ required: true, message: "请输入接口地址" }],
                token: [{ required: true, message: "请输入认证Token" }],
                maxToken: [{ required: true, message: "请输入消息长度" }],
                priority: [{ required: true, message: "请输入优先级" }],
                models: [{ required: true, message: "请输入支持模型" }],
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
