//@flow
import React from 'react';
import {NavLink} from 'react-router-dom';
import {connect} from 'react-redux';
import {Container, Menu} from 'semantic-ui-react';

import type {Settings} from '../types';

type Props = {
  settings: Settings,
};

const NavMenuContainer = (props: Props) => {
  const {settings} = props;
  return (
    <Menu pointing secondary>
      <Container>
        <NavLink exact to="/" className="item" activeClassName="active">
          Home
        </NavLink>
        <NavLink to="/playground" className="item" activeClassName="active">
          Playground
        </NavLink>
        {!settings.isReadOnly && (
          <NavLink to="/questions" className="item" activeClassName="active">
            Q&A
          </NavLink>
        )}
      </Container>
    </Menu>
  );
};

const mapStateToProps = state => {
  return {
    settings: state.settings,
  };
};

const NavMenu = connect(mapStateToProps)(NavMenuContainer);

export default NavMenu;
