<template>
    <div>
        <a-page-header style="border: 1px solid rgb(235, 237, 240); padding: 8px; margin: 0;" title="函数定义"
            @back="back" />
        <a-card>
            <a-form :model="model" :rules="vRules" :label-col="{ style: { width: '90px' } }">
                <a-form-item label="函数名称" name="funcName">
                    <a-input v-model:value="model.funcName" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="绑定角色" name="roleId">
                    <a-select v-model:value="model.roleId" :options="systemDatas" :allowClear="true" :readonly="isReadOnly"/>
                </a-form-item>
                <a-form-item label="函数内容" name="funcContent">
                    <a-textarea v-model:value="model.funcContent" :rows="14" :readonly="isReadOnly" />
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
import { Add, Edit, Get } from '../../../wailsjs/go/biz/FunctionBiz.js';
import { GetList as PromptList } from '../../../wailsjs/go/biz/PromptBiz.js';

export default {
    mounted() {
        this.getSystemDatas();
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
            systemDatas: [],
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
        getSystemDatas() {
            PromptList().then(res => {
                this.systemDatas = res.data;
            });
        },
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
