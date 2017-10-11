<?php
class Scheduler_M extends MY_Model
{

		function __construct ()
	{
		parent::__construct();
	}

	public function create(){

		$today = date("Y/m/d");
		$EmployeeId = 104;
		$HigherManagerId = 1;
		$Status = 1;
		$BidType = 1;
		$companyId = 1;
		$userId = 104;
		if(isset($_POST))
		{

		$master = "insert into scheduler_master (SchedulerName,EmployeeId,HigherManagerId,Status,BidType,CompanyId) values('".$_POST['tripName']."',".$EmployeeId.",".$HigherManagerId.",".$Status.",'".$BidType."','".$companyId."')";
		$this->db->query($master);

		$masterId = $this->db->insert_id();

		//end

		$cityArr = $_POST['cities'];
		$count = count($cityArr);

		for($i=0;$i<$count;$i++){

		// get citymasterid by cityname 

		$this->db->select('cityMasterId,stateMasterId');
		$this->db->from('citymaster');
		$this->db->where('cityName',$cityArr[$i]['title']);
		$citymasterId=$this->db->get();	

		foreach ($citymasterId->result_array() as $result) {	

		//end

		$sql = "insert into scheduler_details (CompanyId,SchedulerMasterId,TripName,CityMasterId,stateMasterId,Checkin,Checkout,PlaceToVisite,Purpose,Status,Created_date,Updated_date,CreaterId,UpdaterId) values('".$companyId."','".$masterId."','".$cityArr[$i]['title']."','".$result['cityMasterId']."','".$result['stateMasterId']."','".$cityArr[$i]['start']."','".$cityArr[$i]['end']."','".$cityArr[$i]['place']."','".$cityArr[$i]['purpose']."','".$Status."','".$today."','".$today."','".$userId."','".$userId."')";
		
		$this->db->query($sql);
			

		}
	}
		if($this->db->affected_rows()>0){

			echo true;
		}
		else{

			echo false;
		}

	}

	}

public function getschedule(){

		$userName['employeeName'] = $this->session->userdata['displayName'];
		
		$userid= $this->session->userdata['id'];
		
		// get benefitBundleId

		$this->db->select('BenefitBundleID');
		$this->db->from('policy_benefit_bundle');
		$this->db->where('CompanyID',$userid);
		$benefitBundleId=$this->db->get()->row()->BenefitBundleID;
		
		// get BenefitTypeID

		$this->db->select('BenefitTypeID');
		$this->db->from('benefit_bundle_type_mapping');
		$this->db->where('BenefitBundleID',$benefitBundleId);
		$benefitType=$this->db->get()->row()->BenefitTypeID;

		//get schedule budget price
		
		$query = $this->db->query("SELECT MaxAmount from benefit_type_allowance_mapping where BenefitBundleTypeMappingID IN (SELECT BenefitBundleTypeMappingID from benefit_bundle_type_mapping where BenefitBundleID = '".$benefitBundleId."' and BenefitTypeID = '".$benefitType."') AND CityCatID IN (SELECT CityCatID from city_mapping WHERE CityID = '1' AND CityCatID IN (SELECT CityCatID from city_category where CompanyID = '104'))");
		

		foreach($query->result_array() as $getet)
		{	
			$details = $getet;
		}

		// get sheduler master details 

			$this->db->select('SchedulerMasterId as tripId,SchedulerName as tripName,CompanyId as companyId,EmployeeId as employeeId,Status as status');
			$this->db->from('scheduler_master');
			$this->db->where('EmployeeId','104');
			$this->db->where('Status','1');
			$schedMaster = $this->db->get();
			
			foreach($schedMaster->result_array() as $details)
		{	
			$value = $details;
			
			//print_r($details);
		}
		$MasterDetails = $value;
		$schedulerMaster = array_merge($MasterDetails,$userName); 
		$schedMasterIds = array();
		
			//get schedulerMasterId

			$this->db->select('SchedulerMasterId');
			$this->db->from('scheduler_master');
			$this->db->where('EmployeeId','104');
			$this->db->where('Status','1');
			$schedMasterId = $this->db->get();
			$fullTrip= array();
			foreach ($schedMasterId->result_array() as $resultIds ) {
				
				$this->db->select('scheduler_details.SchedulerId as id,scheduler_details.Checkin as start,scheduler_details.Checkout as end,scheduler_details.Purpose as purpose,scheduler_details.PlaceToVisite as place,scheduler_details.TripName as title,scheduler_details.Status as status,citymaster.cityName as city,statemaster.stateName as state');
		$this->db->from('scheduler_details');
		$this->db->join('citymaster','citymaster.cityMasterId=scheduler_details.CityMasterId');
	    $this->db->join('statemaster','statemaster.stateMasterId=scheduler_details.stateMasterId');
		$this->db->where_in('SchedulerMasterId',$resultIds);
		$this->db->where('Status','1');
		$scheduler = $this->db->get();
		

		$quotes = array(array("quoteId" => "1",
         "budget" => "2000",
         "hotelName" => "blaha")
		);

		foreach($scheduler->result_array() as $results)
		{		
				// echo "<pre>";
			 // print_r($results);
		$results['start'] = date("Y-m-d", $results['start']);
		$results['end'] = date("Y-m-d", $results['end']);

			$results['quotes']=(array)$quotes;
			array_push($fullTrip,$results);
		}
				
			}
			
		
		
		$schedulerMaster['cities']= $fullTrip;
	
		$schedulerDetails = json_encode($schedulerMaster);
		 echo $schedulerDetails ; 
				
	}

	public function tripEdit(){

		

		$today = date("Y/m/d");
		$EmployeeId = 104;
		$HigherManagerId = 1;
		$Status = 1;
		$BidType = 1;
		$companyId = 2;
		$userId = 104;
		$citymasterId = 1;
		
		
		if(isset($_POST))
		{
			$this->db->select('cityMasterId,stateMasterId');
		$this->db->from('citymaster');
		$this->db->where('cityName',$_POST['title']);
		$citymasterId=$this->db->get();	

		foreach ($citymasterId->result_array() as $result) {
		
			$data = array(
	        		'CompanyId' => $companyId,
	        		'TripName' => $_POST['title'],
	        		'CityMasterId' => $result['cityMasterId'],
	        		'stateMasterId' => $result['stateMasterId'],
	        		'Checkin' => $_POST['start'],
	        		'Checkout' => $_POST['end'],
	        		'PlaceToVisite' => $_POST['place'],
	        		'Purpose' => $_POST['purpose'],
	        		'Status' => $Status,
	        		'Updated_date' => $today,
	        		'UpdaterId' => $userId,
				);

			$this->db->where('SchedulerId', $_POST['id']);
			$this->db->update('scheduler_details', $data);
			
			

		}
	}
		if($this->db->affected_rows()>0){

			echo true;
		}
		else{

			echo false;
		}
			
	}


	public function tripDelete(){

		if(isset($_POST))
		{
				
			$this->db->where('SchedulerId', $_POST['id']);
			$this->db->delete('scheduler_details');
		
		}
		if($this->db->affected_rows()>0){

			echo true;
		}
		else{

			echo false;
		}


	}

	
	public function deleteTrip($param){

		// $_1 = $this->db->query("DELETE FROM scheduler_details WHERE SchedulerMasterId=".$param['tripId']);
		// $_2 = $this->db->query("DELETE FROM scheduler_master WHERE SchedulerMasterId=".$param['tripId']);

			$data = array(
	        		
	        		'Status' => 0,
	        		
				);

			$this->db->where('SchedulerMasterId', $param['tripId']);
			$this->db->update('scheduler_master', $data);

			$this->db->where('SchedulerMasterId', $param['tripId']);
			$this->db->update('scheduler_details', $data);

		if($this->db->affected_rows()>0){

			return true;
		}
		else{
			return false;
		}
	}

	public function getSingle($param){
		$id = $param['id'];
		$res = $this->db->query("SELECT SchedulerMasterId as id FROM scheduler_details WHERE SchedulerId = ".$id);
		$sid = $res->result_array();
		$res2 = $this->db->query("SELECT SchedulerMasterId as tripId,SchedulerName as tripName,status FROM scheduler_master WHERE SchedulerMasterId =".$sid[0]['id']);
		$data = $res2->result_array();
		echo json_encode($data[0]);

	}

	public function cityAndState(){
		$_1 = $this->db->query("SELECT cityName as label,cityMasterId as value FROM citymaster WHERE countryMasterId=101");
		$_2 = $this->db->query("SELECT stateName as label,stateMasterId as value FROM statemaster WHERE countryMasterId=101");
		$res = $_1->result_array();
		$res1 = $_2->result_array();
		echo json_encode(array('state' => $res1,'city' => $res));
	}

	public function getNearByHotel($param){
		$id = $param['id'];
		$res = $this->db->query("SELECT cityMasterId as id FROM scheduler_details WHERE SchedulerId = ".$id);
		$sid = $res->result_array();
		$q = $this->db->query("SELECT hotelname, hotelsmasterid, hclCode, '0' as price FROM hotelsmaster 
		
		as hm, hotelClassifications as hcl where 
		
		hcl.hotelClassificationsId=hm.hotelClassificationsId and citymasterid=".$sid[0]['id']);
		$res1 = $q->result_array();
		echo json_encode($res1);
		

	}

	
	// create template

	public function createTemplate(){
		
		print_r($_POST);

		$today = date("Y/m/d");
		$EmployeeId = 104;
		$HigherManagerId = 1;
		$Status = 1;
		$BidType = 1;
		$companyId = 1;
		$userId = 104;
		if(isset($_POST))
		{

		$master = "insert into scheduler_template (TemplateName,EmployeeId,HigherManagerId,CompanyId) values('".$_POST['tripName']."',".$EmployeeId.",".$HigherManagerId.",'".$companyId."')";
		$this->db->query($master);

		$masterId = $this->db->insert_id();

		//end

		$cityArr = $_POST['cities'];
		$count = count($cityArr);

		for($i=0;$i<$count;$i++){

		// get citymasterid by cityname TemplateMasterId

		$this->db->select('cityMasterId,stateMasterId');
		$this->db->from('citymaster');
		$this->db->where('cityName',$cityArr[$i]['title']);
		$citymasterId=$this->db->get();	

		foreach ($citymasterId->result_array() as $result) {
			

		//end

		$sql = "insert into scheduler_template_details (CompanyId,TemplateMasterId,TripName,CityMasterId,stateMasterId,Checkin,Checkout,PlaceToVisite,Purpose,CreatedDate,updatedDate,CreaterId) values('".$companyId."','".$masterId."','".$cityArr[$i]['title']."','".$result['cityMasterId']."','".$result['stateMasterId']."','".$cityArr[$i]['start']."','".$cityArr[$i]['end']."','".$cityArr[$i]['place']."','".$cityArr[$i]['purpose']."','".$today."','".$today."','".$userId."')";
		
		$this->db->query($sql);	

		}
	}
		if($this->db->affected_rows()>0){

			echo true;
		}
		else{

			echo false;
		}

	}

	}

	public function getTemplates(){

		//$userName['employeeName'] = $this->session->userdata['displayName'];
		$userid= 104;

			$this->db->select('TemplateName as label,TemplateMasterId as value');
			$this->db->from('scheduler_template');
			$this->db->where('EmployeeId',$userid);
			$schedTemplate = $this->db->get();
			
			$templ = array();
			foreach($schedTemplate->result_array() as $details)
		{	

			
			array_push($templ,$details);
				
		}
		echo json_encode(array("templates" => $templ));

	}

	public function getTemplatesList(){

			$this->db->select('scheduler_template_details.TemplateDetailId as id,scheduler_template_details.Checkin as start,scheduler_template_details.Checkout as end,scheduler_template_details.Purpose as purpose,scheduler_template_details.PlaceToVisite as place,scheduler_template_details.TripName as title,citymaster.cityName as state,');
				$this->db->from('scheduler_template_details');
				$this->db->join('citymaster','citymaster.cityMasterId=scheduler_template_details.CityMasterId');
				$this->db->where_in('TemplateMasterId','4');
				$scheduler = $this->db->get();
					$arr = array();
		foreach($scheduler->result_array() as $results)
		{		

			$results['start'] = date("Y-m-d", $results['start']);
			$results['end'] = date("Y-m-d", $results['end']);
				//array("cities"=>$arr)
			array_push($arr,$results);
		}
		echo json_encode(array("cities" => $arr));
	}
}