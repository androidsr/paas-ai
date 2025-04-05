<template>
  <div class="node-panel">
    <div class="node-item" @mousedown="openSelection">
      <a-space direction="vertical" :size="3">
        <FlagTwoTone />
        <span class="node-label">选区</span>
      </a-space>
    </div>

    <div class="node-item" @mousedown="addStartNode">
      <a-space direction="vertical" :size="3">
        <MinusCircleTwoTone />
        <span class="node-label">开始</span>
      </a-space>
    </div>
    <div class="node-item" @mousedown="addLlmTask">
      <a-space direction="vertical" :size="3">
        <MessageTwoTone />
        <span class="node-label">AI问答</span>
      </a-space>
    </div>
    <div class="node-item" @mousedown="addVectorTask">
      <a-space direction="vertical" :size="3">
        <FileWordTwoTone />
        <span class="node-label">知识库</span>
      </a-space>
    </div>
    <div class="node-item" @mousedown="addDatabaseTask">
      <a-space direction="vertical" :size="3">
        <DatabaseTwoTone />
        <span class="node-label">数据库</span>
      </a-space>
    </div>
    <div class="node-item" @mousedown="addFunctionTask">
      <a-space direction="vertical" :size="3">
        <CalculatorTwoTone />
        <span class="node-label">工具</span>
      </a-space>
    </div>
    <div class="node-item" @mousedown="addMcpTask">
      <a-space direction="vertical" :size="3">
        <ControlOutlined style="color: #1890FF;font-size: 18px;" />
        <span class="node-label">MCP</span>
      </a-space>
    </div>
    <div class="node-item" @mousedown="addScriptTask">
      <a-space direction="vertical" :size="3">
        <CodeTwoTone />
        <span class="node-label">代码块</span>
      </a-space>
    </div>
    <div class="node-item" @mousedown="addHttpTask">
      <a-space direction="vertical" :size="3">
        <RocketTwoTone />
        <span class="node-label">HTTP</span>
      </a-space>
    </div>
    <div class="node-item" @mousedown="addEndNode">
      <a-space direction="vertical" :size="3">
        <StopTwoTone />
        <span class="node-label">结束</span>
      </a-space>
    </div>
  </div>
</template>

<script>
import DatabaseNode from "./nodes/ai/DatabaseNode.js";
import EndNode from "./nodes/ai/EndNode.js";
import FunctionNode from "./nodes/ai/FunctionNode.js";
import HttpNode from "./nodes/ai/HttpNode.js";
import LLMNode from "./nodes/ai/LLMNode.js";
import ScriptNode from "./nodes/ai/ScriptNode.js";
import StartNode from "./nodes/ai/StartNode.js";
import VectorNode from "./nodes/ai/VectorNode.js";
import McpNode from "./nodes/ai/McpNode.js";

export default {
  name: "NodePanel",
  props: {
    lf: Object,
  },

  mounted() {
    let lf = this.$props.lf;
    if (lf) {
      lf.on("selection:selected", () => {
        lf.updateEditConfig({ stopMoveGraph: false });
      });
    }
    this.registerNode();
  },
  data() {
    return {
    };
  },
  methods: {
    registerNode() {
      let lf = this.$props.lf;
      lf.register(LLMNode);
      lf.register(VectorNode);
      lf.register(StartNode);
      lf.register(EndNode);
      lf.register(FunctionNode);
      lf.register(HttpNode);
      lf.register(DatabaseNode);
      lf.register(ScriptNode);
      lf.register(McpNode);
    },
    openSelection(event) {
      this.lf.openSelectionSelect();
      this.lf.once("selection:selected", () => {
        this.lf.closeSelectionSelect();
      });
      event.preventDefault();
    },
    addStartNode(event) {
      this.lf.dnd.startDrag({ type: "StartNode", text: "开始" });
      event.preventDefault();
    },
    addLlmTask(event) {
      this.lf.dnd.startDrag({ type: "LLMNode", text: "LLM对话" });
      event.preventDefault();
    },
    addVectorTask(event) {
      this.lf.dnd.startDrag({ type: "VectorNode", text: "知识库" });
      event.preventDefault();
    },
    addFunctionTask(event) {
      this.lf.dnd.startDrag({ type: "FunctionNode", text: "工具调用" });
      event.preventDefault();
    },
    addMcpTask(event) {
      this.lf.dnd.startDrag({ type: "McpNode", text: "MCP" });
      event.preventDefault();
    },
    addScriptTask(event) {
      this.lf.dnd.startDrag({ type: "ScriptNode", text: "代码块" });
      event.preventDefault();
    },
    addDatabaseTask(event) {
      this.lf.dnd.startDrag({ type: "DatabaseNode", text: "DB调用" });
      event.preventDefault();
    },
    addHttpTask(event) {
      this.lf.dnd.startDrag({ type: "HttpNode", text: "HTTP调用" });
      event.preventDefault();
    },
    addEndNode(event) {
      this.lf.dnd.startDrag({ type: "EndNode", text: "结束流程" });
      event.preventDefault();
    },
    /**工作流 */
    addStartEvent(event) {
      this.lf.dnd.startDrag({ type: "StartEvent", text: "开始" });
      event.preventDefault();
    },
    addUserTask(event) {
      this.lf.dnd.startDrag({ type: "UserTask", text: "用户任务" });
      event.preventDefault();
    },
    addExclusiveGateway(event) {
      this.lf.dnd.startDrag({ type: "ExclusiveGateway", text: "" });
      event.preventDefault();
    },
    addParallelGateway(event) {
      this.lf.dnd.startDrag({ type: "ParallelGateway", text: "" });
      event.preventDefault();
    },
    addInclusiveGateway(event) {
      this.lf.dnd.startDrag({ type: "InclusiveGateway", text: "" });
      event.preventDefault();
    },
    addEndEvent(event) {
      this.lf.dnd.startDrag({ type: "EndEvent", text: "结束" });
      event.preventDefault();
    },
  },
};
</script>

<style>
.node-panel {
  position: absolute;
  top: 10px;
  left: 10px;
  width: 64px;
  background-color: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  text-align: center;
  z-index: 101;
}

.node-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 8px;
  cursor: grab;
  /* 添加可拖动的鼠标样式 */
  user-select: none;
  /* 防止文字选中影响拖动 */
  transition: transform 0.2s, background-color 0.3s;

}

.node-label {
  font-size: 12px;
}

.node-item:active {
  cursor: grabbing;
  /* 拖拽过程中鼠标样式 */
}

.node-item:hover {
  background-color: #f0f2f5;
  transform: scale(1.05);
}
</style>
