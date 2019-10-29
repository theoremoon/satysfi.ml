<template>
    <div>
        <form action="">
            <div><input type="text" placeholder="username" v-model="username"></div>
            <div><input type="password" placeholder="password" v-model="password"></div>
        </form>
        <div>
            <div>
                <router-link to="/">Index</router-link>
            </div>
            <div>
                <router-link to="/register">Register</router-link>
            </div>
        </div>
    </div>
</template>

<script>
import Vue from 'vue'

export default Vue.extend({
    data() {
        return {
            username: "",
            password: "",
        }
    },
})

</script>
