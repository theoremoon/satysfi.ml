<template>
    <div class="wrapper">
        <template v-if="sidebar">
            <div class="sidebar-background" @click="sidebar = false"></div>
            <div class="sidebar" @click="sidebar = false">
                <FileTree :tree="files"></FileTree>
            </div>
        </template>
        <div class="menu">
            <button @click="sidebar = true">FILES</button>
            <button>NEW PROJECT</button>
            <button>SAVE</button>
            <button>COMPILE</button>

        </div>
        <div class="main">
            <div ref="editor" id="editor" class="editor">
            </div>
            <div class="viewer">
                <div class="selector">
                    <div @click="tabIndex = 0">PDF</div>
                    <div @click="tabIndex = 1">stdout</div>
                    <div @click="tabIndex = 2">stderror</div>
                </div>
                <div class="content">
                    <template v-if="tabIndex == 0">
                        <embed :src="'data:application/pdf;base64,' + pdf" type="application/pdf" v-if="pdf">
                        <div v-else>No Output</div>
                    </template>
                    <div v-if="tabIndex==1">stdout</div>
                    <div v-if="tabIndex==2">stderr</div>
                </div>
            </div>
        </div>
    </div>
</template>


<script>
import Vue from 'vue'
import axios from 'axios'
import * as monaco from 'monaco-editor'
import FileTree from './FileTree.vue'


const api_data = '{"name":"/","path":"/","childdirs":[{"name":"assets","path":"/assets","childdirs":[],"children":[{"name":"satysfi-logo.jpg","path":"/assets/satysfi-logo.jpg"}]}],"children":[{"name":"demo.saty","path":"/demo.saty"},{"name":"local.satyh","path":"/local.satyh"}]}'
const id = "28fc3a2c3a66faba"

const getFiles = async function(id) {
    return await axios.get("/api/" + id + "/list").then(r => r.data)
}

export default Vue.extend({
    components: {
        FileTree,
    },
    data() {
        return {
            tabIndex: 0,
            pdf: '',
            sidebar: false,
            files: [],
        }
    },
    async mounted() {
        monaco.editor.create(this.$refs.editor, {
            language: 'satysfi',
            automaticLayout: true,
            theme: 'vs'
        })
        this.files = await getFiles(id);
        console.log(this.files)
    },
})
</script>


<style>
.monaco-editor {
    width: 100%;
    height: 100%;
}
</style>
<style scoped lang="scss">
.selector {
    width: 100%;
    display: flex;
    flex-direction: row;
    border-bottom: 1px solid #ccc;

    div {
        border: 1px solid #ccc;
        border-bottom: none;
        border-radius: 2px 2px 0 0 ;
        padding: 0.1em 0.5em;
        margin: 0 0.5em;

        &:hover {
            background-color: rgba(0,0,255, 0.1);
            cursor: pointer;
        }
    }
}
.wrapper {
    width: 100%;
    height: 100%;
}
.menu {
    width: 100%;
    height: 50px;
    background-color: rgba(0,128,128,0.2);
}
.main {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: row;
}
.editor {
    width: 50%;
    height: 100%;
    resize: horizontal;
    overflow: auto;
}

.viewer {
    flex: 1;
    height: 100%;
}
.sidebar-background {
    position: absolute;
    width: 100%;
    height: 100%;
    background-color: rgba(0,0,0,0.7);
    z-index: 999;
}
.sidebar {
    position: fixed;
    width: 40%;
    max-width: 400px;
    min-width: 100px;
    height: 100%;
    background-color: rgba(32,128,32,1);
    z-index:1000;
}
</style>
