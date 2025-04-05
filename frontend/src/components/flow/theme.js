// theme.js
const theme = {
  // 节点的基础样式
  baseNode: {
    fill: "#F0F5FF",             // 浅灰蓝色背景，避免纯白的单调
    stroke: "#1890FF",           // Ant Design 蓝色边框
    strokeDasharray: "none",     // 实线边框
    shadow: {                    // 提升质感的阴影效果
      color: "#D9E6FF",          // 柔和的蓝色阴影
      blur: 8,                   // 增加模糊半径
      offsetX: 2,                // 水平偏移
      offsetY: 2                 // 垂直偏移
    }
  },

  // 矩形节点的样式
  rect: {
    fill: "#F0F5FF",             // 浅灰蓝色背景
    strokeDasharray: "none",     // 实线边框
    radius: 8,                   // 圆角
    stroke: "#1890FF",           // Ant Design 蓝色边框
    shadow: {                    // 阴影效果
      color: "#D9E6FF",          // 阴影颜色
      blur: 8,                   // 模糊半径
      offsetX: 2,                // 水平偏移
      offsetY: 2                 // 垂直偏移
    }
  },

  // 圆形节点的样式
  circle: {
    r: 22,                       // 圆形半径增大
    fill: "#E6F7FF",             // 浅蓝色背景
    stroke: "#1890FF",           // 蓝色边框
    strokeWidth: 1,              // 边框宽度减小
    strokeDasharray: "none",     // 实线边框
    shadow: {                    // 阴影效果
      color: "#D9E6FF",          // 阴影颜色
      blur: 6,                   // 模糊半径
      offsetX: 1,                // 水平偏移
      offsetY: 1                 // 垂直偏移
    }
  },

  // 边的样式（调整为更细的边框）
  edge: {
    strokeWidth: 1,              // 边的宽度减小
    stroke: "#1890FF",           // Ant Design 蓝色边框
    color: "#1890FF",            // 连接线的默认颜色
    strokeDasharray: "none",     // 实线连接线
  },

  // 文字样式
  text: {
    fontSize: 14,                 // 字体大小
    fontWeight: "500",            // 半粗体
    color: "#333333",             // 深灰色文字
    lineHeight: 1.5,              // 行高
    textAlign: 'left',            // 水平靠左对齐
    justifyContent: 'center',     // 如果在一个 flex 容器内，也可以使用此来水平居中
    display: 'flex',              // 使用 flex 布局
    alignItems: 'center',         // 垂直居中
    height: '100%'                // 确保容器高度占满父级
  },

  // 箭头样式
  arrow: {
    offset: 8,                   // 箭头长度
    fill: "none",                // 无填充
    stroke: "#1890FF"            // 蓝色边框
  },

  // 折线连接线样式（调整为更细的边框）
  polyline: {
    stroke: "#1890FF",           // 默认深蓝色边框
    strokeWidth: 2,              // 连接线宽度减小
    strokeDasharray: "none",     // 实线连接线
  },

  // 贝塞尔曲线连接线样式（调整为更细的边框）
  bezier: {
    stroke: "#1890FF",           // 默认深蓝色边框
    strokeWidth: 2,              // 连接线宽度减小
    strokeDasharray: "none",     // 实线连接线
  }
};

// 导出主题配置
export default theme;
