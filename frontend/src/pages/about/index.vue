<template>
    <a-layout style="min-height: 92vh;background-color: white;">
        <a-layout-content>
            <a-row justify="center" align="middle" style="height: 100%; padding: 20px;">
                <a-col :span="18">
                    <a-card title="免责声明" bordered>
                        <a-typography-paragraph>
                            <strong>使用目的：</strong>
                            本软件仅供合法用途使用。用户应确保使用本软件时遵守所在地区的法律法规，禁止用于任何违法活动。
                        </a-typography-paragraph>
                        <a-typography-paragraph>
                            <strong>免责声明：</strong>
                            本软件不对因使用该软件而引起的任何法律责任、损失、或其他不良后果负责。使用本软件的风险由用户自行承担。
                        </a-typography-paragraph>
                        <a-typography-paragraph>
                            <strong>合规性声明：</strong>
                            用户应确保其使用本软件的行为不侵犯任何第三方的合法权益，包括但不限于侵犯版权、专利、商标或其他知识产权，不违反任何法律法规。
                        </a-typography-paragraph>
                        <a-typography-paragraph>
                            <strong>作者信息：</strong>
                            <div style="display: flex; justify-content: space-between; align-items: center;">
                                <div>
                                    <strong>作者：</strong> 那个谁
                                    <br />
                                    <strong>微信号：</strong> scrqrr
                                    <br />
                                    <strong>邮箱：</strong> srandroid@163.com
                                </div>
                                <a-image :src="imageSrc" width="150px" height="150px" alt="扫码增加微信好友"
                                    :preview="false" />
                            </div>
                        </a-typography-paragraph>
                        <!-- <a-space>
                            <a-button type="primary" @click="yesClick" :disabled="$store.isAbout">同意</a-button>
                            <a-button @click="noClick" :disabled="$store.isAbout">不同意</a-button>
                        </a-space> -->
                    </a-card>
                </a-col>
            </a-row>
        </a-layout-content>
    </a-layout>
</template>

<script>
import { Close, GetConfig, SetConfig } from '../../../wailsjs/go/main/App'

export default {
    created() {
        import('@/assets/images/2code.png').then(image => {
            this.imageSrc = image.default;
        });
    },
    data() {
        return {
            imageSrc: null,
        }
    },
    methods: {
        yesClick() {
            SetConfig({ id: "1", legalStatement: "1", content: "{}" }).then(res => {
                if (res.code == 200) {
                    GetConfig().then((res) => {
                        this.$message.success('您已同意声明');
                        if (res.code != 200) {
                            return;
                        }
                        if (res.data.legalStatement == 1) {
                            this.$store.setIsAbout(true);
                            this.$router.push("/taskList");
                        }
                    })
                }
            })
        },
        noClick() {
            //SetConfig({ id: "1", legalStatement: "0" });
            this.$message.error('您不已同意声明');
            Close() // 关闭窗口，并返回 false
        },
    },
};
</script>

<style scoped>
/* 页面样式 */
.a-layout-content {
    padding: 50px 0;
}

.a-card {
    padding: 20px;
    text-align: left;
}
</style>
