<template>
  <div class="node-panel">
    <div class="node-item" @mousedown="openSelection">
      <div class="node-item-icon bpmn-selection"></div>
      <span class="node-label">选区</span>
    </div>
    <div class="node-item" @mousedown="addStartNode">
      <div class="node-item-icon node-start"></div>
      <span class="node-label">开始</span>
    </div>
    <div class="node-item" @mousedown="addLlmTask">
      <div class="node-item-icon node-llm"></div>
      <span class="node-label">AI问答</span>
    </div>
    <div class="node-item" @mousedown="addVectorTask">
      <div class="node-item-icon node-db"></div>
      <span class="node-label">知识库</span>
    </div>
    <div class="node-item" @mousedown="addDatabaseTask">
      <div class="node-item-icon node-database"></div>
      <span class="node-label">DB</span>
    </div>
    <div class="node-item" @mousedown="addFunctionTask">
      <div class="node-item-icon node-func"></div>
      <span class="node-label">函数</span>
    </div>
    <div class="node-item" @mousedown="addHttpTask">
      <div class="node-item-icon node-http"></div>
      <span class="node-label">HTTP</span>
    </div>
    <div class="node-item" @mousedown="addBugTask">
      <div class="node-item-icon node-bugs"></div>
      <span class="node-label">爬虫</span>
    </div>
    <div class="node-item" @mousedown="addEndNode">
      <div class="node-item-icon node-end"></div>
      <span class="node-label">结束</span>
    </div>
  </div>
</template>

<script>
import LogicFlow from '@logicflow/core';
import LLMNode from "@/components/flow/nodes/LLMNode.js";
import VectorNode from "@/components/flow/nodes/VectorNode.js";
import StartNode from "@/components/flow/nodes/StartNode.js";
import EndNode from "@/components/flow/nodes/EndNode.js";
import FunctionNode from "@/components/flow/nodes/FunctionNode.js";
import HttpNode from "@/components/flow/nodes/HttpNode.js";
import DatabaseNode from "@/components/flow/nodes/DatabaseNode.js";
import BugNode from "@/components/flow/nodes/BugNode.js";

export default {
  name: "NodePanel",
  data() {
    return {
      isMousemove: false,
    }
  },
  props: {
    lf: Object,
  },
  mounted() {
    //选区框选使用的
    let lf = this.$props.lf;
    lf &&
      lf.on("selection:selected", () => {
        lf.updateEditConfig({
          stopMoveGraph: false,
        });
      });
    this.registerNode();
  },
  methods: {
    registerNode() {
      let that = this;
      let lf = that.$props.lf;
      lf.register(LLMNode);
      lf.register(VectorNode);
      lf.register(StartNode);
      lf.register(EndNode);
      lf.register(FunctionNode);
      lf.register(HttpNode);
      lf.register(BugNode);
      lf.register(DatabaseNode);
    },
    openSelection(event) {
      this.$props.lf.openSelectionSelect();
      this.$props.lf.once('selection:selected', () => {
        this.$props.lf.closeSelectionSelect();
      });
      /* (this.$props.lf).updateEditConfig({
        stopMoveGraph: true
      }); */
      event.preventDefault();
    },
    addStartNode(event) {
      (this.$props.lf).dnd.startDrag({
        type: "StartNode",
        text: "开始",
      });
      event.preventDefault();
    },
    addLlmTask(event) {
      (this.$props.lf).dnd.startDrag({
        type: "LLMNode",
        text: "LLM对话",
        properties: {},
      });
      event.preventDefault();
    },
    addVectorTask(event) {
      (this.$props.lf).dnd.startDrag({
        type: "VectorNode",
        text: "知识库",
      });
      event.preventDefault();
    },
    addFunctionTask(event) {
      (this.$props.lf).dnd.startDrag({
        type: "FunctionNode",
        text: "函数调用",
      });
      event.preventDefault();
    },
    addDatabaseTask(event) {
      (this.$props.lf).dnd.startDrag({
        type: "DatabaseNode",
        text: "DB调用",
      });
      event.preventDefault();
    },
    addHttpTask(event) {
      (this.$props.lf).dnd.startDrag({
        type: "HttpNode",
        text: "HTTP调用",
      });
      event.preventDefault();
    },
    addBugTask(event) {
      (this.$props.lf).dnd.startDrag({
        type: "BugNode",
        text: "网页爬虫",
      });
      event.preventDefault();
    },

    addEndNode(event) {
      (this.$props.lf).dnd.startDrag({
        type: "EndNode",
        text: "结束",
      });
      event.preventDefault();
    },
  },
};

</script>
<style>
.node-panel {
  margin: 10px;
  position: absolute;
  top: 40px;
  left: 0px;
  width: 60px;
  padding: 5px;
  background-color: white;
  box-shadow: 0 0 10px 1px rgb(228, 224, 219);
  border-radius: 4px;
  text-align: center;
  z-index: 101;
}

.node-item {
  margin-bottom: 10px;
}

.node-item-icon {
  width: 30px;
  height: 30px;
  margin-left: 10px;
  background-size: cover;
}

.node-label {
  font-size: 12px;
  margin-top: 5px;
  user-select: none;
}

.bpmn-selection {
  background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiI+PHBhdGggZmlsbD0ibm9uZSIgc3Ryb2tlPSJjdXJyZW50Q29sb3IiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCIgc3Ryb2tlLXdpZHRoPSIzMiIgZD0ibTM1LjQgODcuMTJsMTY4LjY1IDE5Ni40NEExNi4wNyAxNi4wNyAwIDAgMSAyMDggMjk0djExOS4zMmE3LjkzIDcuOTMgMCAwIDAgNS4zOSA3LjU5bDgwLjE1IDI2LjY3QTcuOTQgNy45NCAwIDAgMCAzMDQgNDQwVjI5NGExNi4wNyAxNi4wNyAwIDAgMSA0LTEwLjQ0TDQ3Ni42IDg3LjEyQTE0IDE0IDAgMCAwIDQ2NiA2NEg0Ni4wNUExNCAxNCAwIDAgMCAzNS40IDg3LjEyIi8+PC9zdmc+) center center no-repeat;
  cursor: pointer;
  filter: invert(45%) sepia(83%) saturate(365%) hue-rotate(123deg) brightness(98%) contrast(86%);
}

.node-start {
  background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiI+PHBhdGggZmlsbD0ibm9uZSIgc3Ryb2tlPSJjdXJyZW50Q29sb3IiIHN0cm9rZS1taXRlcmxpbWl0PSIxMCIgc3Ryb2tlLXdpZHRoPSIzMiIgZD0iTTQ0OCAyNTZjMC0xMDYtODYtMTkyLTE5Mi0xOTJTNjQgMTUwIDY0IDI1NnM4NiAxOTIgMTkyIDE5MnMxOTItODYgMTkyLTE5MloiLz48cGF0aCBmaWxsPSJjdXJyZW50Q29sb3IiIGQ9Ik0zMjAgMTc2YTE2IDE2IDAgMCAwLTE2IDE2djUzbC0xMTEuNjgtNjcuNDRhMTAuNzggMTAuNzggMCAwIDAtMTYuMzIgOS4zMXYxMzguMjZhMTAuNzggMTAuNzggMCAwIDAgMTYuMzIgOS4zMUwzMDQgMjY3djUzYTE2IDE2IDAgMCAwIDMyIDBWMTkyYTE2IDE2IDAgMCAwLTE2LTE2Ii8+PC9zdmc+) center center no-repeat;
  cursor: pointer;
  filter: invert(45%) sepia(83%) saturate(365%) hue-rotate(123deg) brightness(98%) contrast(86%);
}

.node-end {
  background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiI+PHBhdGggZmlsbD0ibm9uZSIgc3Ryb2tlPSJjdXJyZW50Q29sb3IiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCIgc3Ryb2tlLXdpZHRoPSIzMiIgZD0iTTQxNiA0NDhIOTZhMzIuMDkgMzIuMDkgMCAwIDEtMzItMzJWOTZhMzIuMDkgMzIuMDkgMCAwIDEgMzItMzJoMzIwYTMyLjA5IDMyLjA5IDAgMCAxIDMyIDMydjMyMGEzMi4wOSAzMi4wOSAwIDAgMS0zMiAzMiIvPjwvc3ZnPg==) center center no-repeat;
  cursor: pointer;
  filter: invert(45%) sepia(83%) saturate(365%) hue-rotate(123deg) brightness(98%) contrast(86%);
}

.node-llm {
  background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIwLjk5ZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgMjU2IDI2MCI+PHBhdGggZD0iTTIzOS4xODQgMTA2LjIwM2E2NC43MTYgNjQuNzE2IDAgMCAwLTUuNTc2LTUzLjEwM0MyMTkuNDUyIDI4LjQ1OSAxOTEgMTUuNzg0IDE2My4yMTMgMjEuNzRBNjUuNTg2IDY1LjU4NiAwIDAgMCA1Mi4wOTYgNDUuMjJhNjQuNzE2IDY0LjcxNiAwIDAgMC00My4yMyAzMS4zNmMtMTQuMzEgMjQuNjAyLTExLjA2MSA1NS42MzQgOC4wMzMgNzYuNzRhNjQuNjY1IDY0LjY2NSAwIDAgMCA1LjUyNSA1My4xMDJjMTQuMTc0IDI0LjY1IDQyLjY0NCAzNy4zMjQgNzAuNDQ2IDMxLjM2YTY0LjcyIDY0LjcyIDAgMCAwIDQ4Ljc1NCAyMS43NDRjMjguNDgxLjAyNSA1My43MTQtMTguMzYxIDYyLjQxNC00NS40ODFhNjQuNzY3IDY0Ljc2NyAwIDAgMCA0My4yMjktMzEuMzZjMTQuMTM3LTI0LjU1OCAxMC44NzUtNTUuNDIzLTguMDgzLTc2LjQ4M20tOTcuNTYgMTM2LjMzOGE0OC4zOTcgNDguMzk3IDAgMCAxLTMxLjEwNS0xMS4yNTVsMS41MzUtLjg3bDUxLjY3LTI5LjgyNWE4LjU5NSA4LjU5NSAwIDAgMCA0LjI0Ny03LjM2N3YtNzIuODVsMjEuODQ1IDEyLjYzNmMuMjE4LjExMS4zNy4zMi40MDkuNTYzdjYwLjM2N2MtLjA1NiAyNi44MTgtMjEuNzgzIDQ4LjU0NS00OC42MDEgNDguNjAxTTM3LjE1OCAxOTcuOTNhNDguMzQ1IDQ4LjM0NSAwIDAgMS01Ljc4MS0zMi41ODlsMS41MzQuOTIxbDUxLjcyMiAyOS44MjZhOC4zMzkgOC4zMzkgMCAwIDAgOC40NDEgMGw2My4xODEtMzYuNDI1djI1LjIyMWEuODcuODcgMCAwIDEtLjM1OC42NjVsLTUyLjMzNSAzMC4xODRjLTIzLjI1NyAxMy4zOTgtNTIuOTcgNS40MzEtNjYuNDA0LTE3LjgwM00yMy41NDkgODUuMzhhNDguNDk5IDQ4LjQ5OSAwIDAgMSAyNS41OC0yMS4zMzN2NjEuMzlhOC4yODggOC4yODggMCAwIDAgNC4xOTUgNy4zMTZsNjIuODc0IDM2LjI3MmwtMjEuODQ1IDEyLjYzNmEuODE5LjgxOSAwIDAgMS0uNzY3IDBMNDEuMzUzIDE1MS41M2MtMjMuMjExLTEzLjQ1NC0zMS4xNzEtNDMuMTQ0LTE3LjgwNC02Ni40MDV6bTE3OS40NjYgNDEuNjk1bC02My4wOC0zNi42M0wxNjEuNzMgNzcuODZhLjgxOS44MTkgMCAwIDEgLjc2OCAwbDUyLjIzMyAzMC4xODRhNDguNiA0OC42IDAgMCAxLTcuMzE2IDg3LjYzNXYtNjEuMzkxYTguNTQ0IDguNTQ0IDAgMCAwLTQuNC03LjIxM20yMS43NDItMzIuNjlsLTEuNTM1LS45MjJsLTUxLjYxOS0zMC4wODFhOC4zOSA4LjM5IDAgMCAwLTguNDkyIDBMOTkuOTggOTkuODA4Vjc0LjU4N2EuNzE2LjcxNiAwIDAgMSAuMzA3LS42NjVsNTIuMjMzLTMwLjEzM2E0OC42NTIgNDguNjUyIDAgMCAxIDcyLjIzNiA1MC4zOTF6TTg4LjA2MSAxMzkuMDk3bC0yMS44NDUtMTIuNTg1YS44Ny44NyAwIDAgMS0uNDEtLjYxNFY2NS42ODVhNDguNjUyIDQ4LjY1MiAwIDAgMSA3OS43NTctMzcuMzQ2bC0xLjUzNS44N2wtNTEuNjcgMjkuODI1YTguNTk1IDguNTk1IDAgMCAwLTQuMjQ2IDcuMzY3em0xMS44NjgtMjUuNThMMTI4LjA2NyA5Ny4zbDI4LjE4OCAxNi4yMTh2MzIuNDM0bC0yOC4wODYgMTYuMjE4bC0yOC4xODgtMTYuMjE4eiIvPjwvc3ZnPg==) center center no-repeat;
  cursor: pointer;
  filter: invert(45%) sepia(83%) saturate(365%) hue-rotate(123deg) brightness(98%) contrast(86%);
}

.node-db {
  background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiI+PHBhdGggZmlsbD0iY3VycmVudENvbG9yIiBkPSJtNDc5LjY2IDI2OC43bC0zMi0xNTEuODFDNDQxLjQ4IDgzLjc3IDQxNy42OCA2NCAzODQgNjRIMTI4Yy0xNi44IDAtMzEgNC42OS00Mi4xIDEzLjk0cy0xOC4zNyAyMi4zMS0yMS41OCAzOC44OWwtMzIgMTUxLjg3QTE2LjcgMTYuNyAwIDAgMCAzMiAyNzJ2MTEyYTY0IDY0IDAgMCAwIDY0IDY0aDMyMGE2NCA2NCAwIDAgMCA2NC02NFYyNzJhMTYuNyAxNi43IDAgMCAwLS4zNC0zLjNtLTM4NC0xNDUuNHYtLjI4YzMuNTUtMTguNDMgMTMuODEtMjcgMzIuMjktMjdIMzg0YzE4LjYxIDAgMjguODcgOC41NSAzMi4yNyAyNi45MWMwIC4xMy4wNS4yNi4wNy4zOWwyNi45MyAxMjcuODhhNCA0IDAgMCAxLTMuOTIgNC44MkgzMjBhMTUuOTIgMTUuOTIgMCAwIDAtMTYgMTUuODJhNDggNDggMCAxIDEtOTYgMEExNS45MiAxNS45MiAwIDAgMCAxOTIgMjU2SDcyLjY1YTQgNCAwIDAgMS0zLjkyLTQuODJaIi8+PC9zdmc+) center center no-repeat;
  cursor: pointer;
  filter: invert(45%) sepia(83%) saturate(365%) hue-rotate(123deg) brightness(98%) contrast(86%);
}

.node-http {
  background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgMjQgMjQiPjxwYXRoIGZpbGw9IiMwODkxYjIiIGZpbGwtcnVsZT0iZXZlbm9kZCIgZD0iTTE3LjMwMyA2LjY5N2E3LjUgNy41IDAgMCAwLTEyLjQ3IDMuMDg4QTUuMDAyIDUuMDAyIDAgMCAwIDYuNSAxOS41aDEyYTQgNCAwIDAgMCAuOTktNy44NzZhNy40NzQgNy40NzQgMCAwIDAtMi4xODctNC45MjciIGNsaXAtcnVsZT0iZXZlbm9kZCIvPjwvc3ZnPg==) center center no-repeat;
  cursor: pointer;
  filter: invert(45%) sepia(83%) saturate(365%) hue-rotate(123deg) brightness(98%) contrast(86%);
}

.node-database {
  background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgMjQgMjQiPjxwYXRoIGZpbGw9IiMwODkxYjIiIGZpbGwtcnVsZT0iZXZlbm9kZCIgZD0iTTEuMjUgOUE2Ljc1IDYuNzUgMCAwIDEgOCAyLjI1aDhBNi43NSA2Ljc1IDAgMCAxIDIyLjc1IDl2NkE2Ljc1IDYuNzUgMCAwIDEgMTYgMjEuNzVIOEE2Ljc1IDYuNzUgMCAwIDEgMS4yNSAxNXptMTYuMDgtLjA5NmMtLjI1LS4xNC0uNTg4LS4xNTQtLjk3My0uMTU0SDE0LjI1djIuNWgyLjEyOGMuMzc3IDAgLjcwNi0uMDE3Ljk1Mi0uMTU0Yy4xNTgtLjA4OS40Mi0uMzAyLjQyLTEuMDk2cy0uMjYyLTEuMDA3LS40Mi0xLjA5Nk0xOC41ODIgMTJjLjQyMy0uNDUuNjY4LTEuMTEyLjY2OC0yYzAtMS4yMDYtLjQ1Mi0xLjk5My0xLjE4Ny0yLjQwNGMtLjYyLS4zNDgtMS4zMjUtLjM0Ni0xLjY3LS4zNDZIMTMuNWEuNzUuNzUgMCAwIDAtLjc1Ljc1djhjMCAuNDE0LjMzNi43NS43NS43NWgyLjg5NGMuMzQ0IDAgMS4wNDkuMDAyIDEuNjY5LS4zNDZjLjczNS0uNDExIDEuMTg3LTEuMTk4IDEuMTg3LTIuNDA0YzAtLjg4OC0uMjQ1LTEuNTQ5LS42NjgtMm0tMi4yMDMuNzVIMTQuMjV2Mi41aDIuMTA3Yy4zODUgMCAuNzIzLS4wMTQuOTczLS4xNTRjLjE1OC0uMDg5LjQyLS4zMDIuNDItMS4wOTZzLS4yNjItMS4wMDctLjQyLTEuMDk1Yy0uMjQ1LS4xMzgtLjU3NS0uMTU0LS45NTEtLjE1NW0tOC4wMjItNGMuMzg1IDAgLjcyMy4wMTQuOTczLjE1NGMuMTU4LjA4OS40Mi4zMDIuNDIgMS4wOTZ2NGMwIC43OTQtLjI2MiAxLjAwNy0uNDIgMS4wOTZjLS4yNS4xNC0uNTg4LjE1NC0uOTczLjE1NEg2LjI1di02LjV6TTExLjI1IDEwYzAtMS4yMDYtLjQ1Mi0xLjk5My0xLjE4Ny0yLjQwNGMtLjYyLS4zNDgtMS4zMjUtLjM0Ni0xLjY3LS4zNDZINS41YS43NS43NSAwIDAgMC0uNzUuNzV2OGMwIC40MTQuMzM2Ljc1Ljc1Ljc1aDIuODk0Yy4zNDQgMCAxLjA0OS4wMDIgMS42NjktLjM0NmMuNzM1LS40MTEgMS4xODctMS4xOTggMS4xODctMi40MDR6IiBjbGlwLXJ1bGU9ImV2ZW5vZGQiLz48L3N2Zz4=) center center no-repeat;
  cursor: pointer;
  filter: invert(45%) sepia(83%) saturate(365%) hue-rotate(123deg) brightness(98%) contrast(86%);
}

.node-func {
  background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiI+PHBhdGggZD0iTTM1MiA0NDIuMmMtLjMtMi4yLTItMy45LTQuMi00LjNsLTIyLjMtMS45Yy0xMS44LTMuMS0yMC41LTE2LjItMjIuMy0yOC4zTDMwMiA0MDBoLTkybC0xLjIgNy42Yy0xLjkgMTIuMS0xMC41IDI1LjItMjIuMyAyOC4zbC0yMi4zIDEuOWMtMi4xLjUtMy45IDIuMi00LjIgNC4zLS40IDMuMSAyIDUuOCA1LjEgNS44aDE4MS44YzMgLjEgNS41LTIuNiA1LjEtNS43eiIgZmlsbD0iI2ZmZmZmZiIvPjxwYXRoIGQ9Ik00NzIuOSA3MWMtNC41LTQuNS0xMC43LTctMTctN0g1Ni4yYy02LjQgMC0xMi41IDIuNS0xNyA3UzMyIDgxLjcgMzIgODh2MjcyYzAgNi40IDIuNyAxMi41IDcuMiAxN3MxMC41IDcgMTYuOCA3aDQwMGM2LjQgMCAxMi4zLTIuNSAxNi44LTdzNy4yLTEwLjcgNy4yLTE3Vjg4Yy4xLTYuMy0yLjYtMTIuNS03LjEtMTd6TTI1NiAzNjBjLTQuNCAwLTgtMy42LTgtOHMzLjYtOCA4LTggOCAzLjYgOCA4LTMuNiA4LTggOHptMjAzLjktNDBINTIuMmMtMi4yIDAtNC0xLjgtNC00VjkyYzAtNi42IDUuNC0xMiAxMi0xMmgzOTEuN2M2LjYgMCAxMiA1LjQgMTIgMTJ2MjI0YzAgMi4yLTEuOCA0LTQgNHoiIGZpbGw9IiNmZmZmZmYiLz48L3N2Zz4=) center center no-repeat;
  cursor: pointer;
  filter: invert(45%) sepia(83%) saturate(365%) hue-rotate(123deg) brightness(98%) contrast(86%);
}

.node-bugs {
  background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxZW0iIGhlaWdodD0iMWVtIiB2aWV3Qm94PSIwIDAgMjQgMjQiPjxwYXRoIGZpbGw9IiMwODkxYjIiIGQ9Ik0xMiAyMXEtMS42IDAtMi44NjMtLjgyNVQ3LjA3NiAxOGwtMS41MjUuODc1cS0uMzUuMi0uNzUuMDc1dC0uNi0uNDc1dC0uMS0uNzV0LjQ1LS42bDEuNzI1LTFxLS4wNzUtLjI3NS0uMTI1LS41NjNUNi4wNSAxNUg0cS0uNDI1IDAtLjcxMi0uMjg4VDMgMTR0LjI4OC0uNzEyVDQgMTNoMi4wNXEuMDUtLjMuMS0uNTg3dC4xMjUtLjU2M2wtMS43MjUtMXEtLjM1LS4yLS40NS0uNnQuMS0uNzV0LjYxMy0uNDYydC43NjIuMDg3TDcuMDUgMTBxLjItLjM1LjQ2My0uNjg3VDguMDUgOC43UTggOC41MjUgOCA4LjM1VjhxMC0uNi4xNzUtMS4xNXQuNDc1LTEuMDI1bC0uOTUtLjk1cS0uMjc1LS4yNzUtLjI4OC0uN1Q3LjcgMy40NXEuMjc1LS4zLjY4OC0uMjg3dC43MTIuMjg3bDEuMDUgMXEuNDI1LS4yMjUuODg4LS4zMzdUMTIgNHQuOTc1LjEyNXQuOS4zNUwxNC45IDMuNDVxLjMtLjMuNy0uMjg4dC43LjMxM3EuMjc1LjMuMjg4Ljd0LS4yODguN2wtLjk1Ljk1cS4zLjQ3NS40NjMgMS4wMjVUMTUuOTc1IDh2LjMzOHEwIC4xNjItLjA1LjMzN3EuMjc1LjI3NS41MzguNjI1dC40NjIuN2wxLjUyNS0uODc1cS4zNS0uMi43NS0uMDg4dC42LjQ2M3QuMDg4Ljc2M3QtLjQ2My42MTJsLTEuNzI1Ljk3NXEuMDc1LjI3NS4xMzguNTYzdC4xMTIuNTg3SDIwcS40MjUgMCAuNzEzLjI4OFQyMSAxNHQtLjI4OC43MTNUMjAgMTVoLTIuMDVxLS4wNS4zLS4xLjU4OHQtLjEyNS41NjJsMS43MjUgMXEuMzUuMi40NS42MTN0LS4xLjc2MnQtLjYuNDV0LS43NS0uMUwxNi45MjUgMThxLS44IDEuMzUtMi4wNjMgMi4xNzVUMTIgMjFNMTAuMSA3LjM1cS40MjUtLjE3NS45MTMtLjI2MlQxMiA3dC45NjMuMDc1dC44ODcuMjVxLS4yLS41NzUtLjctLjk1VDEyIDZ0LTEuMTc1LjM4OHQtLjcyNS45NjJNMTIgMTlxMS44MjUgMCAyLjkxMy0xLjUyNVQxNiAxNHEwLTEuNzUtMS4wMTItMy4zNzVUMTIgOXEtMS45NSAwLTIuOTc1IDEuNjEzVDggMTRxMCAxLjk1IDEuMDg4IDMuNDc1VDEyIDE5bTAtMnEtLjQyNSAwLS43MTItLjI4OFQxMSAxNnYtNHEwLS40MjUuMjg4LS43MTJUMTIgMTF0LjcxMy4yODhUMTMgMTJ2NHEwIC40MjUtLjI4OC43MTNUMTIgMTciLz48L3N2Zz4=) center center no-repeat;
  cursor: pointer;
  filter: invert(45%) sepia(83%) saturate(365%) hue-rotate(123deg) brightness(98%) contrast(86%);
}
</style>