/*
 * Copyright 2015, Robert Bieber
 *
 * This file is part of mixer.
 *
 * mixer is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mixer is distributed in the hope that it will be useful,
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with mixer.  If not, see <http://www.gnu.org/licenses/>.
 */

import React from 'react';
import Qajax from 'qajax';

import Intro from '../views/Intro.js';
import Composer from '../views/Composer.js';
import Submitted from '../views/Submitted.js';

export default class IndexPage extends React.Component {
	constructor(props, context) {
		super(props, context);
		this.state = {
			token: null,
			playlists: null,
			userID: null,
		};
	}

	onLogin(data) {
		setTimeout(this.refreshToken.bind(this), data.expires_in * 1000);
		this.setState(
			{token: data.token},
			() => Qajax({
				url: this.props.playlistsURI,
				params: {token: this.state.token},
			})
				.then(Qajax.filterSuccess)
				.then(Qajax.toJSON)
				.then(payload => this.setState(payload))
		)
	}

	onSubmit(data) {
		this.setState(
			{submitted: true},
			() => console.log(data)
		);
	}

	refreshToken() {
		console.log('Refreshing login tokens');
	}

	render() {
		var view = (
			<Intro
				loginURI={this.props.loginURI}
				onLogin={this.onLogin.bind(this)}
			/>
		);

		if (this.state.token !== null) {
			view = (
				<Composer
					{...this.state}
					onSubmit={this.onSubmit.bind(this)}
				/>
			);
		}

		if (this.state.submitted) {
			view = <Submitted />;
		}

		return (
			<div className="container">
				{view}
			</div>
		);
	}
}
IndexPage.propTypes = {
	loginURI: React.PropTypes.string.isRequired,
	playlistsURI: React.PropTypes.string.isRequired,
};
