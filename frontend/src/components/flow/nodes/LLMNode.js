import { RectResize } from "@logicflow/extension";

import { h } from 'preact';

// 提供节点
class LLMNode extends RectResize.view {
  getResizeShape() {
    // 获取XxxNodeModel中定义的形状属性
    const { model } = this.props;
    const { x, y, width, height, radius, properties } = model;
    // 获取XxxNodeModel中定义的样式属性
    const style = model.getNodeStyle();

    return h('g', {}, [
      h('rect', {
        ...style,
        x: x - width / 2,
        y: y - height / 2,
        width,
        height,
        rx: radius,
        ry: radius,
      }),
      h('svg', {
        x: x - width / 2 + 3,
        y: y - height / 2 + 10,
        width: 70,
        height: 70,
        viewBox: "0 0 1024 1024",
      }, [
        h('path', {
          fill: style.stroke,
          d: "M239.184 106.203a64.716 64.716 0 0 0-5.576-53.103C219.452 28.459 191 15.784 163.213 21.74A65.586 65.586 0 0 0 52.096 45.22a64.716 64.716 0 0 0-43.23 31.36c-14.31 24.602-11.061 55.634 8.033 76.74a64.665 64.665 0 0 0 5.525 53.102c14.174 24.65 42.644 37.324 70.446 31.36a64.72 64.72 0 0 0 48.754 21.744c28.481.025 53.714-18.361 62.414-45.481a64.767 64.767 0 0 0 43.229-31.36c14.137-24.558 10.875-55.423-8.083-76.483m-97.56 136.338a48.397 48.397 0 0 1-31.105-11.255l1.535-.87l51.67-29.825a8.595 8.595 0 0 0 4.247-7.367v-72.85l21.845 12.636c.218.111.37.32.409.563v60.367c-.056 26.818-21.783 48.545-48.601 48.601M37.158 197.93a48.345 48.345 0 0 1-5.781-32.589l1.534.921l51.722 29.826a8.339 8.339 0 0 0 8.441 0l63.181-36.425v25.221a.87.87 0 0 1-.358.665l-52.335 30.184c-23.257 13.398-52.97 5.431-66.404-17.803M23.549 85.38a48.499 48.499 0 0 1 25.58-21.333v61.39a8.288 8.288 0 0 0 4.195 7.316l62.874 36.272l-21.845 12.636a.819.819 0 0 1-.767 0L41.353 151.53c-23.211-13.454-31.171-43.144-17.804-66.405zm179.466 41.695l-63.08-36.63L161.73 77.86a.819.819 0 0 1 .768 0l52.233 30.184a48.6 48.6 0 0 1-7.316 87.635v-61.391a8.544 8.544 0 0 0-4.4-7.213m21.742-32.69l-1.535-.922l-51.619-30.081a8.39 8.39 0 0 0-8.492 0L99.98 99.808V74.587a.716.716 0 0 1 .307-.665l52.233-30.133a48.652 48.652 0 0 1 72.236 50.391zM88.061 139.097l-21.845-12.585a.87.87 0 0 1-.41-.614V65.685a48.652 48.652 0 0 1 79.757-37.346l-1.535.87l-51.67 29.825a8.595 8.595 0 0 0-4.246 7.367zm11.868-25.58L128.067 97.3l28.188 16.218v32.434l-28.086 16.218l-28.188-16.218z",
        })
      ])

    ]);
  }
}

// 提供节点的属性
class LLMNodeModel extends RectResize.model {
  constructor(data, graphModel) {
    super(data, graphModel);
  }

  initNodeData(data) {
    console.log(data)
    super.initNodeData(data);
    this.width = 200;
    this.height = 36;
    this.radius = 2;
    this.text.draggable = true;
  }
  getNodeStyle() {
    const style = super.getNodeStyle();
    return style;
  }
}

export default {
  type: "LLMNode",
  view: LLMNode,
  model: LLMNodeModel
}