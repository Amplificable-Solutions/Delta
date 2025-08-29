pragma solidity 0.8.10;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract TestERC20A is ERC20 {
	event __DeltaSendToAccount(address recipient, uint256 amount);

	constructor() public ERC20("Bitcoin MAX", "MAX") {
		_mint(msg.sender, 100000000000000000000000000);
	}

	function test_native_transfer(uint amount) public {
		emit __DeltaSendToAccount(msg.sender, amount);
	}
}
