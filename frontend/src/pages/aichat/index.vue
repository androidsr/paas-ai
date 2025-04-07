<template>
  <div style="background-color: white;">
    <a-row :gutter="16">
      <a-col :span="24">
        <a-row>
          <div :style="'width:100%;height:100%;border: 1px solid #e9e9e9;'">
            <Message :messages="$store.chats.result" ref="msgView" />
            <!-- 
            <md-preview :modelValue="result" :codeFoldable="false" :preview-theme="'github'" :theme="'light'">
            </md-preview> -->
          </div>
        </a-row>
        <a-row>
          <a-col>
            <a-space style="padding-top:10px;" size="middle">
              <a-popover placement="bottomRight">
                <template #content>
                  <a-image-preview-group v-if="chatType == '6'">
                    <a-image :width="120" v-for="item in uploadImageUrls" :src="item" v-if="!!this.uploadFileId" />
                  </a-image-preview-group>
                </template>
                <a-button size="small">
                  <CloudUploadOutlined @click="fileChange"
                    title="多模态模型：支持图片、文件(csv、txt、pdf、xlsx、docx、pptx、md、html)上传解读对话" />
                </a-button>
              </a-popover>

              <a-button size="small">
                <CopyOutlined ref="copyBtn" />
              </a-button>

              <a-button size="small" @click="clean" title="清空对话">
                <ClearOutlined />
              </a-button>
              <a-popover placement="rightBottom" trigger="click">
                <template #content>
                  <div>
                    <a-form style="max-height: 75vh;overflow: auto;" :label-col="{ style: { width: '80px' } }">
                      <a-form-item label="AI平台">
                        <a-select v-model:value="channelId" :options="ahChannels" @select="getChannelItem"
                          @change="settingOk" />
                      </a-form-item>
                      <a-form-item label="选择模型">
                        <a-select v-model:value="model" :options="models" @change="settingOk" />
                      </a-form-item>
                      <a-form-item label="对话模式">
                        <a-select v-model:value="chatType" :options="chatTypeData" @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="函数调用" v-if="chatType === '4' || chatType === '9'">
                        <a-select v-model:value="funcCall" :options="funcs" :allowClear="true" @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="选择角色">
                        <a-select v-model:value="systemId" :options="systemDatas" :allowClear="true"
                          @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="系统提示词">
                        <a-textarea v-model:value="system" :disabled="!!systemId" @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="历史数量">
                        <a-input-number v-model:value="cacheLimit" style="width: 100%" @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="流程调用" v-if="chatType === '5'">
                        <a-select v-model:value="flowId" :options="flows" :allowClear="true" @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="知识库集合" v-if="chatType === '2'">
                        <a-select v-model:value="collectionId" :options="collectionList" @select="collectionChange"
                          :allowClear="true" @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="知识库文档" v-if="chatType === '2'">
                        <a-select v-model:value="filename" :options="filenameList" :allowClear="true"
                          @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="最大条数" v-if="chatType === '2'">
                        <a-input-number v-model:value="pageSize" style="width: 100%" @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="相似比例" v-if="chatType === '2'">
                        <a-slider v-model:value="similarityScore" :min="0" :max="10" @change="settingOk" />
                      </a-form-item>
                      <a-form-item label="温度控制" title="控制输出结果的随机性：值越低更严谨、值越高创建性强。">
                        <a-slider v-model:value="temperature" :min="0" :max="10" @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="采样(P)" title="控制文本的随机性和多样性。">
                        <a-slider v-model:value="topP" :min="0" :max="10" @change="settingOk" />
                      </a-form-item>

                      <a-form-item label="采样(K)" title="token保留数量，选出排名前N个值">
                        <a-slider v-model:value="topK" :min="0" :max="20" @change="settingOk" />
                      </a-form-item>
                      <a-form-item label="输出设置">
                        <a-checkbox v-model:checked="stream" @change="settingOk">流式输出</a-checkbox>
                        <a-space :size="10" />
                        <a-checkbox v-model:checked="cnAnswer" @change="settingOk">强制中文</a-checkbox>
                        <a-space :size="10" />
                        <a-checkbox v-model:checked="jsonFormat" @change="settingOk">格式化</a-checkbox>
                      </a-form-item>
                    </a-form>
                  </div>
                </template>
                <a-button size="small">
                  <SettingOutlined />
                </a-button>
              </a-popover>
            </a-space>
          </a-col>
        </a-row>
        <a-row style="margin-top: 10px;">
          <a-col :flex="1">
            <a-textarea type="text" placeholder="Enter发送" :rows="3" v-model:value="message"
              @keydown.enter.prevent="sendMessage" :loading="loadingSend" style="resize: none; " />
          </a-col>
        </a-row>
      </a-col>
    </a-row>

    <input v-show="false" multiple :disabled="isLoading" ref="inputFile" type="file" @change="addFile($event)" />
  </div>
</template>
<style>
.container {
  border: 1px solid #ccc;
  overflow: auto;
  padding: 10px;
}
</style>
<script>
import { fetchEventSource } from '@microsoft/fetch-event-source';
import { message } from 'ant-design-vue';
import Clipboard from 'clipboard';
import { Base64 } from 'js-base64';

import { GetList as ChannelList, Get } from '../../../wailsjs/go/biz/AiChannelBiz.js';
import { Clean, GetImage, Upload } from '../../../wailsjs/go/biz/ChatBiz.js';
import { GetList as CollectionList } from '../../../wailsjs/go/biz/CollectionBiz.js';
import { GetList as DocumentList } from '../../../wailsjs/go/biz/DocumentBiz.js';
import { GetList as FunctionList } from '../../../wailsjs/go/biz/FunctionBiz.js';
import { GetList as PromptList } from '../../../wailsjs/go/biz/PromptBiz.js';
import { GetList as FwConfigList } from '../../../wailsjs/go/biz/FwConfigBiz.js';
import { GetConfig } from '../../../wailsjs/go/main/App';
import Message from '../../components/Message.vue';

export default {
  components: {
    Message,
  },
  created() {
    this.getFlowsList();
    this.getSystemDatas();
    this.getFunctionsDatas();
    this.handleResize();
    this.getCollectionList();
    this.getChannelList();
    GetConfig().then(res => {
      if (res.code == 200) {
        this.config = JSON.parse(res.data.content);
        if (!this.channelId) {
          this.getChannelItem(this.config.channelId);
        } else {
          this.getChannelItem(this.channelId);
        }
      } else {
        this.$message.error("获取设置失败！");
      }
    });
  },
  mounted() {
    this.result = this.$store.chats.result;
    this.clipboard = new Clipboard(this.$refs.copyBtn, {
      text: () => {
        let content = "";
        this.$store.chats.result.forEach(item => {
          content += item.message + "\n";
        });
        return content;
      },
    });

    this.clipboard.on('success', () => {
      message.success("复制成功！");
    });
    this.getImage();
    window.addEventListener('resize', this.handleResize)
    this.$store.bus.on("changeChart", (data) => {
      message.info("生成中...");
      const oldSystem = this.system;
      this.system = data;
      const oldChatType = this.chatType;
      this.chatType = "1";
      this.message = "更换图表";
      this.settingOk();
      this.sendMessage();
      this.chatType = oldChatType;
      this.system = oldSystem;
      this.message = "";
      this.settingOk();
    });
  },
  beforeUnmount() {
    this.$store.bus.off('changeChart');
  },
  data() {
    return {
      clipboard: null,
      copySuccess: false,
      ahChannels: [],
      channel: {},
      config: {},
      uploadImageUrls: [],
      tabData: {
        chatConfig: '对话设置',
        modelConfig: '模型设置',
      },
      selected: 'chatConfig',
      chatTypeData: [
        { label: "文本对话", value: "1" },
        { label: "向量库对话", value: "2" },
        { label: "文件对话", value: "3" },
        { label: "工具对话", value: "4" },
        { label: "MCP对话", value: "9" },
        { label: "流程对话", value: "5" },
        { label: "图片对话", value: "6" },
        { label: "音频对话", value: "7" },
        { label: "视频对话", value: "8" },
      ],
      opend: false,
      logsOpend: false,
      file: [],
      selectedFile: "",
      logs: "",
      socket: null,
      width: document.documentElement.clientWidth,
      height: document.documentElement.clientHeight,
      loadingSend: false,
      isLoading: false,
      models: [],
      flows: [],
      pptTemplates: [],
      source: [{ label: '模型对话', value: '1' }, { label: '向量对话', value: '2' }, { label: '文件对话', value: '3' }],
      collectionList: [],
      systemDatas: [],
      funcs: [],
      filenameList: [],
      ctrl: new AbortController(),
      message: "",
      channelId: this.$store.chats.channelId,
      chatType: this.$store.chats.chatType,
      stream: this.$store.chats.stream,
      model: this.$store.chats.model,
      pptTemplate: this.$store.chats.pptTemplate,
      cacheLimit: this.$store.chats.cacheLimit,
      sourceType: this.$store.chats.sourceType,
      system: this.$store.chats.system,
      systemId: this.$store.chats.systemId,
      flowId: this.$store.chats.flowId,
      collectionId: this.$store.chats.collectionId,
      collectionName: this.$store.chats.collectionName,
      cnAnswer: this.$store.chats.cnAnswer,
      similarityScore: this.$store.chats.similarityScore,
      funcCall: this.$store.chats.funcCall,
      temperature: this.$store.chats.temperature,
      jsonFormat: this.$store.chats.jsonFormat,
      topK: this.$store.chats.topK,
      topP: this.$store.chats.topP,
      mode: this.$store.chats.mode,
      filename: this.$store.chats.filename,
      pageSize: this.$store.chats.pageSize,
      uploadFileId: this.$store.chats.uploadFileId,
      result: this.$store.chats.result,
    };
  },
  deactivated() {
    this.clean();
    this.$destroy(true);
    window.removeEventListener('resize', this.handleResize);
  },
  beforeDestroy() {
    if (this.clipboard) {
      this.clipboard.destroy();
    }
  },
  methods: {
    settingOk() {
      //this.showSetting = false;
      var sendData = {
        databaseId: this.databaseId,
        channelId: this.channelId,
        chatType: this.chatType,
        sessionId: this.sessionId,
        system: this.system,
        stream: this.stream,
        systemId: this.systemId,
        model: this.model,
        cacheLimit: this.cacheLimit,
        uploadFileId: this.uploadFileId,
        sourceType: this.sourceType,
        collectionName: this.collectionName,
        cnAnswer: this.cnAnswer,
        collectionId: this.collectionId,
        jsonFormat: this.jsonFormat,
        funcCall: this.funcCall,
        similarityScore: this.similarityScore,
        temperature: this.temperature,
        result: this.result,
        flowId: this.flowId,
        topK: this.topK,
        topP: this.topP,
        filename: this.filename,
        mode: this.mode,
        pageSize: this.pageSize,
        pptTemplate: this.pptTemplate,
        filter: {
          filename: this.filename
        }
      }
      this.$store.setChats(sendData);
    },
    getChannelList() {
      ChannelList().then(res => {
        if (res.code == 200) {
          this.ahChannels = res.data;
        } else {
          this.$message.error("设置保存失败！");
        }
      })
    },
    getChannelItem(channelId) {
      this.channelId = channelId;
      Get(channelId).then(res => {
        if (res.code == 200) {
          this.models = [];
          this.$store.setChannel(res.data);
          console.log(this.$store.channel)

          this.$store.channel.models.split(",").forEach(item => {
            this.models.push({ label: item, value: item });
          });
          this.model = this.models[0].value;
        } else {
          this.$router.push({
            path: "/config"
          });
        }
      });
    },
    getFileTypeValue(fileName) {
      const fileExtension = fileName.split('.').pop().toLowerCase();
      const fileTypes = {
        image: ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'tiff', 'webp', 'svg'],
        audio: ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'wma'],
        video: ['mp4', 'avi', 'mov', 'mkv', 'flv', 'webm', 'wmv', 'mpeg', 'mpg'],
        document: ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt', 'rtf', 'odt', 'csv', 'md', 'html'],
      };

      if (fileTypes.image.includes(fileExtension)) {
        this.chatType = '6'; // 图片对话
      } else if (fileTypes.audio.includes(fileExtension)) {
        this.chatType = '7'; // 音频对话
      } else if (fileTypes.video.includes(fileExtension)) {
        this.chatType = '8'; // 视频对话
      } else if (fileTypes.document.includes(fileExtension)) {
        this.chatType = '3'; // 文件对话
      }
      return this.chatType;
    },
    async getImage() {
      GetImage()
      if (this.uploadImageUrls.length == 0 && !!this.uploadFileId) {
        let ids = this.uploadFileId.split(',');
        for (let i = 0; i < ids.length; i++) {
          let res = await GetImage(ids[i]);
          this.uploadImageUrls[i] = "data:image/png;base64," + res.data;
        }
      }
    },
    sendMessage() {
      if (window.EventSource) {
        if (!this.message) {
          return;
        }
        let newMessage = "#### " + this.message;
        this.$store.addResult({ type: 'user', message: newMessage });
        var sendStr = this.message;
        if (this.cnAnswer) {
          sendStr += ";使用中文回答"
        }
        var sendData = {
          channelId: this.channelId,
          chatType: this.chatType,
          system: this.system,
          stream: this.stream,
          systemId: this.systemId,
          message: sendStr,
          model: this.model,
          cacheLimit: this.cacheLimit,
          uploadFileId: this.uploadFileId,
          sourceType: this.sourceType,
          collectionName: this.collectionName,
          cnAnswer: this.cnAnswer,
          collectionId: this.collectionId,
          jsonFormat: this.jsonFormat,
          funcCall: this.funcCall,
          similarityScore: this.similarityScore,
          temperature: this.temperature,
          result: this.result,
          flowId: this.flowId,
          topK: this.topK,
          topP: this.topP,
          filename: this.filename,
          mode: this.mode,
          pageSize: this.pageSize,
          pptTemplate: this.pptTemplate,
          filter: {
            filename: this.filename
          }
        }
        this.$store.setChats(sendData);
        var m = this;
        var data = Object.assign({ ...sendData }, {
          similarityScore: (this.similarityScore * 0.1).toFixed(2),
          temperature: (this.temperature * 0.1).toFixed(2),
          topK: (this.topK * 0.1).toFixed(2),
          topP: (this.topP * 0.1).toFixed(2),
          result: [],
        });
        var endData;
        var url = "http://localhost:11083";
        switch (this.chatType) {
          case "1":
            url += "/chat/text";
            break;
          case "2":
            url += "/chat/vector";
            break;
          case "3":
            url += "/chat/file";
            break;
          case "4":
            url += "/chat/tools";
            break;
          case "5":
            url += "/chat/flow";
            break;
          case "6":
            url += "/chat/image";
            break;
          case "7":
            message.info("功能开发中...")
            return;
          case "8":
            message.info("功能开发中...")
            return;
          case "9":
            url += "/chat/mcp";
            break;
        }
        fetchEventSource(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          openWhenHidden: true,
          body: JSON.stringify(data),
          signal: m.ctrl.signal,
          onopen: () => {
            m.isLoading = true;
            m.message = "";
            m.$store.addResult({ type: 'ai', message: "" });
          },
          onmessage: (e) => {
            if (this.isLoading) {
              endData = Base64.decode(e.data)
              m.$store.setResult(endData);
            }
          },
          onclose() {
            m.isLoading = false;
            setTimeout(() => {
              m.$refs.msgView.renderAll();
            }, 500);
            message.success("处理完成！");
            m.ctrl.abort();
          },
          onerror(err) {
            console.log(err);
          }
        })
      }
    },
    downFile(filename) {
      location.href = '/chat/download/' + filename
    },
    clearLogs() {
      this.logs = "";
    },
    getSystemDatas() {
      PromptList().then(res => {
        this.systemDatas = res.data;
      });
    },
    getFunctionsDatas() {
      FunctionList().then(res => {
        this.funcs = res.data;
      });
    },
    clean() {
      this.isLoading = false;
      this.message = "";
      this.$store.cleanResult();
      this.result = this.$store.chats.result;
      Clean();
    },
    handleResize() {
      this.width = document.documentElement.clientWidth - 16;
      this.height = document.documentElement.clientHeight - 139;
    },
    sourceChange(data) {
      if (data.value == "2") {
        this.getCollectionList();
      }
    },
    getFlowsList() {
      FwConfigList().then(res => {
        if (res.code == 200) {
          this.flows = res.data;
        } else {
          this.flows = [];
          message.error(res.data.msg);
        }
      })
    },
    getFilenameList(collectionName) {
      this.filenameList = [];
      DocumentList(collectionName).then(res => {
        this.filenameList = res.data;
      });
    },
    getCollectionList() {
      CollectionList().then(res => {
        res.data.forEach(item => {
          this.collectionList.push({ label: item, value: item });
        });
      });
    },
    collectionChange(data) {
      this.filename = "";
      this.collectionName = data;
      this.getFilenameList(data);
    },
    fileChange() {
      this.$refs.inputFile.click();
    },
    addFile(event) {
      this.file = event.target.files;
      const names = [];
      for (let i = 0; i < this.file.length; i++) {
        names.push(this.file[i].name);
      }
      this.selectedFile = names.join(",");
      this.convertFilesToBase64(event);
    },
    // 将文件转换为 Base64 字符串
    async convertFilesToBase64(event) {
      const files = this.file;
      const promises = {};
      // 遍历每个文件，将其转换为 Base64
      for (let i = 0; i < files.length; i++) {
        let f = files[i]
        promises[f.name] = await this.fileToBase64(f)
      }
      this.uploadFile(promises);
    },
    async fileToBase64(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onloadend = () => {
          resolve(reader.result); // 读取文件为 Base64 格式
        };
        reader.onerror = reject;
        reader.readAsDataURL(file); // 将文件读取为 Data URL（Base64 格式）
      });
    },
    uploadFile(base64Files) {
      console.log(base64Files)
      Upload(base64Files).then(res => {
        this.isLoading = false;
        if (res.code == 200) {
          message.info("上传成功");
          this.file = [];
          this.uploadImageUrls = [];
          this.uploadFileId = res.data;
          this.getFileTypeValue(this.uploadFileId);
          if (this.chatType == '6') {
            this.getImage();
          }

        } else {
          message.error(res.msg);
        }
      }).catch(err => {
        this.isLoading = false;
        message.error(err);
      });
    },
  },
  beforeDestroy() {
    if (!!this.uploadImageUrls) {
      this.uploadImageUrls.forEach(uploadImage => {
        URL.revokeObjectURL(uploadImage);
      })
    }
  },
};
</script>

<style scoped>
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

.md-editor {
  height: 78vh;
}
</style>