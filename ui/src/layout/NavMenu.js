import React from "react";
import { NavLink } from "react-router-dom";
import { Container, Menu } from "semantic-ui-react";

const NavMenu = () => {
  return (
    <Menu>
      <Container>
        <NavLink exact to="/" className="item" activeClassName="active">
          Home
        </NavLink>
        <NavLink to="/playground" className="item" activeClassName="active">
          Playground
        </NavLink>
        <NavLink to="/questions" className="item" activeClassName="active">
          Q&A
        </NavLink>
      </Container>
    </Menu>
  );
};

export default NavMenu;
