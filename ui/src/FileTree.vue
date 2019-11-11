<template>
    <div>
        <p>{{ tree.name }}</p>
        <ul>
            <li v-for="file in tree.children" @click="click(file.path)">{{ file.name }}</li>
            <FileTree v-for="dir in tree.childdirs" :tree="dir"></FileTree>
        </ul>
    </div>

</template>

<script>
import Vue from 'vue'
import FileTree from './FileTree.vue'

export default Vue.extend({
    name: 'FileTree',
    props: [
        'tree'
    ],
    methods: {
        click(path) {
            this.$store.dispatch('loadFile', path)
        }
    },
})

</script>

<style lang="scss" scoped>
p,ul {
    margin: 0;
}
li {
    &:hover {
        cursor: pointer;
        background-color: green;
    }
}
</style>
