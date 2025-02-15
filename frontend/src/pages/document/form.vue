<template>
    <div>
        <a-page-header style="border: 1px solid rgb(235, 237, 240); padding: 8px; margin: 0;" title="向量库文档管理"
            @back="back" />
        <a-card>
            <a-form :model="model" :rules="vRules" :label-col="{ style: { width: '90px' } }">
                <a-form-item label="集合名称" name="collectionName">
                    <a-select placeholder="集合名称" v-model:value="model.collectionName" allowClear>
                        <a-select-option v-for="v in collectionList" :value="v">{{ v }}</a-select-option>
                    </a-select>
                </a-form-item>
                <a-form-item label="文件名称" name="filename">
                    <a-input v-model:value="model.filename" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="附加信息" name="metadata">
                    <a-textarea v-model:value="model.metadata" :readonly="isReadOnly" />
                </a-form-item>
                <a-form-item label="文档内容" name="content">
                    <a-textarea v-model:value="model.content" rows="10" :readonly="isReadOnly" />
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
import { Add, Edit, Get } from '../../../wailsjs/go/biz/DocumentBiz.js';
import { GetList } from '../../../wailsjs/go/biz/CollectionBiz.js';

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
        this.getList();
    },
    data() {
        return {
            id: "",
            model: {},
            collectionList: [],
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
        getList() {
            GetList(this.query).then(res => {
                this.collectionList = res.data || [];
            })
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
                    } else {
                        this.$message.error(res.msg);
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
