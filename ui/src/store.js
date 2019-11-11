import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";
import { Base64 } from "js-base64";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    id: null,
    pdf: "",
    stdout: "",
    stderr: "",
    fileTree: "",
    currentFile: null
  },
  getters: {
    id(state) {
      return state.id;
    },
    pdf(state) {
      return state.pdf;
    },
    stdout(state) {
      return state.stdout;
    },
    stderr(state) {
      return state.stderr;
    },
    fileTree(state) {
      return state.fileTree;
    },
    currentFile(state) {
      return state.currentFile;
    }
  },
  mutations: {
    setID(state, id) {
      state.id = id;
    },
    setPDF(state, pdf) {
      state.pdf = pdf;
    },
    setStdoutStderr(state, { stdout, stderr }) {
      state.stdout = stdout;
      state.stderr = stderr;
    },
    setFileTree(state, fileTree) {
      state.fileTree = Vue.util.extend({}, fileTree);
    },
    setCurrentFile(state, { content, name, path }) {
      state.currentFile = {
        content,
        name,
        path
      };
    },
    setContent(state, content) {
      state.currentFile.content = content;
    }
  },
  actions: {
    async newProject(context) {
      const newID = await axios.post("/api/new-project");
      return context.dispatch("loadProject", newID.data.id);
    },
    async loadProject(context, id) {
      context.commit("setID", id);
      return axios.get("/api/" + context.state.id + "/list").then(r => {
        context.commit("setFileTree", r.data);
      });
    },
    async loadFile(context, path) {
      const file = await axios.get("/api/" + context.state.id + "/get", {
        params: {
          path: path
        }
      });
      context.commit("setCurrentFile", file.data);
      return file;
    },
    async newFile(context, path) {
      context.commit("setCurrentFile", {
        content: "",
        path: path,
        name: path
      });
      const promise = context.dispatch("save", "");
      context.dispatch("loadProject", context.state.id);
      return promise;
    },
    async save(context, content) {
      const promise = axios.post("/api/" + context.state.id + "/save", {
        path: context.state.currentFile.path,
        data: Base64.encode(content)
      });
      context.commit("setContent", content);
      return promise;
    },
    async compile(context) {
      return axios
        .post("/api/" + context.state.id + "/compile", {
          path: context.state.currentFile.path
        })
        .then(r => {
          context.commit("setPDF", r.data.pdf);
          context.commit("setStdoutStderr", r.data);
          return r.data;
        });
    }
  }
});
