<template>
    <div class="tag-tree">
        <el-row type="flex" justify="space-between">
            <el-col :span="24" class="el-scrollbar flex-auto" style="overflow: auto">
                <el-input v-model="filterText" placeholder="输入关键字->搜索已展开节点信息" clearable size="small" class="mb5" />

                <el-tree
                    ref="treeRef"
                    :style="{ maxHeight: state.height, height: state.height, overflow: 'auto' }"
                    :highlight-current="true"
                    :indent="10"
                    :load="loadNode"
                    :props="treeProps"
                    lazy
                    node-key="key"
                    :expand-on-click-node="true"
                    :filter-node-method="filterNode"
                    @node-click="treeNodeClick"
                    @node-expand="treeNodeClick"
                    @node-contextmenu="nodeContextmenu"
                >
                    <template #default="{ node, data }">
                        <span>
                            <span v-if="data.type.value == TagTreeNode.TagPath">
                                <tag-info :tag-path="data.label" />
                            </span>

                            <slot v-else :node="node" :data="data" name="prefix"></slot>

                            <span class="ml3">
                                <slot name="label" :data="data"> {{ data.label }}</slot>
                            </span>

                            <slot :node="node" :data="data" name="suffix"></slot>
                        </span>
                    </template>
                </el-tree>
            </el-col>
        </el-row>
        <contextmenu :dropdown="state.dropdown" :items="state.contextmenuItems" ref="contextmenuRef" @currentContextmenuClick="onCurrentContextmenuClick" />
    </div>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, watch, toRefs } from 'vue';
import { TagTreeNode } from './tag';
import TagInfo from './TagInfo.vue';
import Contextmenu from '@/components/contextmenu/index.vue';

const props = defineProps({
    height: {
        type: [Number, String],
        default: 0,
    },
    load: {
        type: Function,
        required: false,
    },
    loadTags: {
        type: Function,
        required: false,
    },
    loadContextmenuItems: {
        type: Function,
        required: false,
    },
});

const treeProps = {
    label: 'name',
    children: 'zones',
    isLeaf: 'isLeaf',
};

const emit = defineEmits(['nodeClick', 'currentContextmenuClick']);
const treeRef: any = ref(null);
const contextmenuRef = ref();

const state = reactive({
    height: 600 as any,
    filterText: '',
    dropdown: {
        x: 0,
        y: 0,
    },
    contextmenuItems: [],
    opend: {},
});
const { filterText } = toRefs(state);

onMounted(async () => {
    if (!props.height) {
        setHeight();
        window.onresize = () => setHeight();
    } else {
        state.height = props.height;
    }
});

const setHeight = () => {
    state.height = window.innerHeight - 157 + 'px';
};

watch(filterText, (val) => {
    treeRef.value?.filter(val);
});

const filterNode = (value: string, data: any) => {
    if (!value) return true;
    return data.label.includes(value);
};

/**
 * 加载树节点
 * @param { Object } node
 * @param { Object } resolve
 */
const loadNode = async (node: any, resolve: any) => {
    if (typeof resolve !== 'function') {
        return;
    }
    let nodes = [];
    try {
        if (node.level == 0 && props.loadTags) {
            nodes = await props.loadTags(node);
        } else if (props.load) {
            nodes = await props.load(node);
        } else {
            nodes = await node.data.loadChildren();
        }
    } catch (e: any) {
        console.error(e);
    }
    return resolve(nodes);
};

const treeNodeClick = (data: any) => {
    emit('nodeClick', data);
    if (data.type.nodeClickFunc) {
        data.type.nodeClickFunc(data);
    }
    // 关闭可能存在的右击菜单
    contextmenuRef.value.closeContextmenu();
};

// 树节点右击事件
const nodeContextmenu = (event: any, data: any) => {
    // 加载当前节点是否需要显示右击菜单
    let items = data.type.contextMenuItems;
    if (!items || items.length == 0) {
        if (props.loadContextmenuItems) {
            items = props.loadContextmenuItems(data);
        }
    }
    if (!items) {
        return;
    }
    state.contextmenuItems = items;
    const { clientX, clientY } = event;
    state.dropdown.x = clientX;
    state.dropdown.y = clientY;
    contextmenuRef.value.openContextmenu(data);
};

const onCurrentContextmenuClick = (clickData: any) => {
    emit('currentContextmenuClick', clickData);
};

const reloadNode = (nodeKey: any) => {
    let node = getNode(nodeKey);
    node.loaded = false;
    node.expand();
};

const getNode = (nodeKey: any) => {
    let node = treeRef.value.getNode(nodeKey);
    if (!node) {
        throw new Error('未找到节点: ' + nodeKey);
    }
    return node;
};

defineExpose({
    reloadNode,
});
</script>

<style lang="scss" scoped>
.tag-tree {
    overflow: 'auto';
    position: relative;

    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
    border: 1px solid var(--el-border-color-light, #ebeef5);

    .el-tree {
        display: inline-block;
        min-width: 100%;
    }
}
</style>
