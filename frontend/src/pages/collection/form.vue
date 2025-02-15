<template>
    <div>
        <a-page-header style="border: 1px solid rgb(235, 237, 240); padding: 8px; margin: 0;" title="向量库集合管理"
            @back="back" />
        <a-card>
            <a-form :model="model" :rules="vRules" :label-col="{ style: { width: '90px' } }">
                <a-form-item label="集合名称" name="collectionName">
                    <a-input v-model:value="collectionName" :readonly="isReadOnly" />
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
import { Add } from '../../../wailsjs/go/biz/CollectionBiz.js';

export default {
    mounted() {
        this.collectionName = this.$route.query.name;
    },
    data() {
        return {
            collectionName: "",
            vRules: {
                collectionName: [{ required: true, message: "请输入集合名称" }],
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
            Add(this.collectionName).then(res => {
                if (res.code == 200) {
                    this.$message.success(res.msg);
                    this.back();
                }
            })
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
