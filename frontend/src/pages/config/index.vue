<template>
    <a-layout style="background-color: white;">
        <a-layout-content>
            <a-row justify="center" align="middle" style="padding: 20px;">
                <a-col :span="18">
                    <a-card title="设置" bordered>
                        <div style="overflow: auto;height: 76vh;" class="no-scrollbar">
                            <a-form :form="form" :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }"
                                @submit.prevent="handleSubmit">
                                <a-form-item label="AI平台">
                                    <a-select v-model:value="form.channelId" :options="ahChannels" />
                                </a-form-item>
                                <a-form-item label="嵌入服务器">
                                    <a-input v-model:value="form.embeddingUrl" />
                                </a-form-item>
                                <a-form-item label="嵌入认证Token">
                                    <a-input v-model:value="form.embeddingToken" />
                                </a-form-item>
                                <a-form-item label="嵌入模型">
                                    <a-input v-model:value="form.embeddingModel" placeholder="bge-m3" />
                                </a-form-item>
                                <a-form-item :wrapper-col="{ span: 3, offset: 20 }">
                                    <a-button type="primary" html-type="submit">保存设置</a-button>
                                </a-form-item>
                            </a-form>
                        </div>
                    </a-card>
                </a-col>
            </a-row>
        </a-layout-content>
    </a-layout>
</template>

<script>
import { GetConfig, SetConfig } from '../../../wailsjs/go/main/App';
import { GetList } from '../../../wailsjs/go/biz/AiChannelBiz.js';

export default {
    mounted() {
        GetConfig().then(res => {
            if (res.code == 200) {
                this.form = JSON.parse(res.data.content);
            } else {
                this.$message.error("获取设置失败！");
            }
        });
        this.getList();
    },
    data() {
        return {
            ahChannels: [],
            channel: {},
            form: {
                channelId: "",
                embeddingUrl: "",
                embeddingToken: "",
                embeddingModel: "bge-m3",
            },
        };
    },
    methods: {
        // 保存设置
        handleSubmit() {
            SetConfig({ id: "1", legalStatement: "1", content: JSON.stringify(this.form) }).then(res => {
                if (res.code == 200) {
                    this.$message.success("设置成功");
                } else {
                    this.$message.error("设置保存失败！");

                }
            })
        },
        getList() {
            GetList().then(res => {
                if (res.code == 200) {
                    this.ahChannels = res.data;
                } else {
                    this.$message.error("设置保存失败！");
                }
            })
        },
    },
};
</script>

<style scoped>
.a-layout-content {
    padding: 50px 0;
}

.a-card {
    padding: 20px;
}
</style>

<style>
.container {
    border: 1px solid #ccc;
    overflow: auto;
    padding: 10px;
}

.no-scrollbar::-webkit-scrollbar {
    display: none;
}

.no-scrollbar {
    -ms-overflow-style: none;
    scrollbar-width: none;
}
</style>