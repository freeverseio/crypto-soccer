import React, { Component } from 'react'; 
import { Image } from 'semantic-ui-react'

class Shop extends Component {
    state = { clicked: false }
    render() {
        const { clicked } = this.state;

        return <Image src={clicked ? "shop/become_sponsor.png" : "shop/shop.png"} onClick={() => this.setState({ clicked: !clicked })} fluid />
    }
}
export default Shop; 