<template>
    <div>
        <a-page-header style="border: 1px solid rgb(235, 237, 240); padding: 8px; margin: 0;" title="数据库连接配置"
            @back="back" />
        <a-card>
            <a-form :model="model" :rules="vRules" :label-col="{ style: { width: '90px' } }">
                <a-form-item label="数据库类型" name="dbType">
                    <a-input v-model:value="model.dbType" :readonly="isReadOnly" placeholder="mysql" />
                </a-form-item>
                <a-form-item label="数据库名称" name="dbname">
                    <a-input v-model:value="model.dbname" rows="25" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="连接地址" name="url">
                    <a-input v-model:value="model.url" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="账户名称" name="username">
                    <a-input v-model:value="model.username" rows="25" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="账户密码" name="password">
                    <a-input v-model:value="model.password" rows="25" :readonly="isReadOnly" />
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
import { Add, Edit, Get } from '../../../wailsjs/go/biz/DbconfigBiz.js';

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
                url: "root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8",
            },
            vRules: {
                dbType: [{ required: true, message: "请输入数据库类型" }],
                dbname: [{ required: true, message: "请输入数据库名称" }],
                url: [{ required: true, message: "请输入连接地址" }],
                username: [{ required: true, message: "请输入账户名称" }],
                password: [{ required: true, message: "请输入账户密码" }],
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
