<template>
    <div class="wrapper">
        <template v-if="sidebar">
            <div class="sidebar-background" @click="sidebar = false"></div>
            <div class="sidebar" @click="sidebar = false">
                <FileTree :tree="files" @path-click="loadFile"></FileTree>
            </div>
        </template>
        <div class="menu">
            <button @click="sidebar = true">FILES</button>
            <button @click="newProject">NEW PROJECT</button>
            <button @click="save">SAVE</button>
            <button @click="compile">COMPILE</button>
            <button @click="newFile">NEW FILE</button>
            <br>
            <div>{{ current_file ? current_file.path : "No file opened" }}</div>
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
                    <pre v-if="tabIndex==1">{{ stdout }}</pre>
                    <pre v-if="tabIndex==2">{{ stderr }}</pre>
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


const getList = async function(id) {
    return await axios.get("/api/" + id + "/list").then(r => r.data)
}
const getFile = async function(id, path) {
    return await axios.get("/api/" + id + "/get", {
        params: {
            path: path,
        }
    })
        .then(r => r.data)
}
const newProject = async function() {
    return await axios.post("/api/new-project").then(r => r.data)
}
const saveFile = async function (id, path, data) {
    return await axios.post("/api/" + id + "/save", {
        path: path,
        data: btoa(data),
    })
}
const compileRequest = async function (id, path) {
    return await axios.post("/api/" + id + "/compile", {
        path: path,
    })
        .then(r => ({
            pdf: r.data.pdf,
            stdout: r.data.stdout,
            stderr: r.data.stderr,
        }))
}

export default Vue.extend({
    components: {
        FileTree,
    },
    data() {
        return {
            id: '',
            tabIndex: 0,
            pdf: '',
            stdout: '',
            stderr: '',
            editor: undefined,
            sidebar: false,
            files: [],
            current_file: undefined, // {name: string, path: string, content: string}
        }
    },
    async mounted() {
        this.editor = monaco.editor.create(this.$refs.editor, {
            language: 'satysfi',
            automaticLayout: true,
            theme: 'satysfier',
            minimap: {
                enabled: false,
            },
        })
        if (this.$route.params.hasOwnProperty('id')) {
            this.loadProject(this.$route.params.id)
        }
    },
    methods: {
        async loadProject(id) {
            this.id = id
            this.files = await getList(this.id);
        },
        async loadFile(path) {
            this.current_file = await getFile(this.id, path)
            this.editor.setValue(this.current_file.content)
        },
        async newProject() {
            const id = (await newProject()).id
            this.$router.push({
                path: `/project/${id}`
            })
            this.loadProject(id)
        },
        async save() {
            saveFile(this.id, this.current_file.path, this.current_file.content)
        },
        async compile() {
            let result = await compileRequest(this.id, this.current_file.path)
            this.stdout = result.stdout
            this.stderr = result.stderr
            this.pdf = result.pdf
        },
        async newFile() {
            let path = window.prompt("Path for new file")
            if (!path) {
                return;
            }
            saveFile(this.id, path, "");
            getList(this.id)
                .then(r => { this.files = r })
            this.loadFile(path)
        },
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
    overflow: overlay;
    padding-bottom: 10px;
    z-index: 888;

}

.viewer {
    flex: 1;
    height: 100%;
}
.sidebar-background {
    position: fixed;
    width: 100%;
    height: 100%;
    background-color: rgba(0,0,0,0.2);
    z-index: 999;
}
.sidebar {
    position: fixed;
    width: 40%;
    max-width: 400px;
    min-width: 100px;
    height: 100%;
    background-color: limegreen;
    z-index:1000;
}

.content {
    width: 100%;
    height: 100%;
    overflow: auto;
}
embed {
    width: 100%;
    height: calc(100% - 10px);
}
</style>
